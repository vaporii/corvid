package srv

import (
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type notifServer server

func (s notifServer) GetCapabilities() (capabilities []string, e *dbus.Error) {
	// log.Print("GetCapabilities called")
	return []string{
		"body",
		"actions",
	}, nil
}

func (s notifServer) GetServerInformation() (name, vendor, version, specVersion string, e *dbus.Error) {
	// log.Print("GetServerInformation called")
	return "corvid", "CartConnoisseur", "0.1.0", "1.2", nil
}

func (s notifServer) CloseNotification(id uint32) (e *dbus.Error) {
	// log.Printf("CloseNotification called: %d", id)
	server(s).close(id, CloseReasonClosed)
	server(s).output()
	return nil
}

func (s notifServer) Notify(appName string, replacesId uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (id uint32, e *dbus.Error) {
	// log.Print("Notify called")
	s.notifications.mutex.Lock()
	defer s.notifications.mutex.Unlock()

	if replacesId == 0 {
		id = s.notifications.nextId
		s.notifications.nextId++
	} else {
		id = replacesId
	}

	actionMap := make(map[string]string)
	for i := 0; i < len(actions)-1; i += 2 {
		actionMap[actions[i]] = actions[i+1]
	}

	hintMap := make(map[string]hint)
	img := ""

	for key, value := range hints {
		if !value.Signature().Empty() {
			if strings.Contains("ybnqiuxtds", string(value.Signature().String()[0])) {
				hintMap[key] = hint{Variant: value}
			} else if key == "image-data" {
				raw := value.Value().([]interface{})

				var i image.Image
				if raw[3].(bool) {
					i = &image.NRGBA{
						Pix:    raw[6].([]uint8),
						Stride: int(raw[2].(int32)),
						Rect:   image.Rect(0, 0, int(raw[0].(int32)), int(raw[1].(int32))),
					}
				} else {
					rgb := raw[6].([]uint8)
					rgba := make([]uint8, len(rgb)/3*4)

					for i := 0; i < len(rgb)-1; i += 3 {
						rgba[i/3*4] = rgb[i]
						rgba[i/3*4+1] = rgb[i+1]
						rgba[i/3*4+2] = rgb[i+2]
						rgba[i/3*4+3] = 0xff
					}

					i = &image.NRGBA{
						Pix:    rgba,
						Stride: int(raw[2].(int32)),
						Rect:   image.Rect(0, 0, int(raw[0].(int32)), int(raw[1].(int32))),
					}
				}
				_ = i

				f, err := os.CreateTemp(os.TempDir(), "corvid-*.png")
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				png.Encode(f, i)

				img = f.Name()
			}
		}
	}

	if expireTimeout == -1 {
		expireTimeout = DEFAULT_EXPIRATION
	}

	notification := notification{
		Id:         id,
		AppName:    appName,
		AppIcon:    appIcon,
		Summary:    summary,
		Body:       body,
		Actions:    actionMap,
		Hints:      hintMap,
		Timestamp:  time.Now().Unix(),
		Expiration: expireTimeout,
		Image:      img,
		timer:      nil,
	}

	if expireTimeout != 0 {
		notification.timer = time.AfterFunc(time.Duration(expireTimeout)*time.Millisecond, func() {
			server(s).close(notification.Id, CloseReasonExpire)
			server(s).output()
		})
	}

	s.notifications.notifications[id] = notification
	server(s).output()

	return id, nil
}
