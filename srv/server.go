package srv

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/godbus/dbus/v5"
)

var DEFAULT_EXPIRATION int
var SORT_DIRECTION int // 1 = newest first, -1 = oldest first

type server struct {
	notifications *notificationStack
	conn          *dbus.Conn
	object        dbus.ObjectPath
	name          string
}

func (s server) close(id uint32, reason closeReason) {
	s.notifications.mutex.Lock()
	defer s.notifications.mutex.Unlock()

	n, ok := s.notifications.notifications[id]
	if !ok {
		return
	}

	if n.timer != nil {
		n.timer.Stop()
	}

	if n.Image != "" {
		os.Remove(n.Image)
	}

	delete(s.notifications.notifications, n.Id)

	err := s.conn.Emit(s.object, s.name+".NotificationClosed", n.Id, reason)
	if err != nil {
		log.Print(err)
	}
}

// TODO: relocate to cmd/corvid
func (s server) output() {
	arr := make([]notification, len(s.notifications.notifications))

	i := 0
	for _, notification := range s.notifications.notifications {
		arr[i] = notification
		i++
	}

	slices.SortFunc(arr, func(a, b notification) int {
		if a.Timestamp > b.Timestamp {
			return SORT_DIRECTION
		} else if a.Timestamp < b.Timestamp {
			return -SORT_DIRECTION
		} else {
			if a.Id > b.Id {
				return SORT_DIRECTION
			} else if a.Id < b.Id {
				return -SORT_DIRECTION
			}
		}

		return 0
	})

	j, err := json.Marshal(arr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(j))
}

func Start(defaultExpiration int, sortDirection int) {
	const NOTIF_DBUS_OBJECT = "/org/freedesktop/Notifications"
	const NOTIF_DBUS_NAME = "org.freedesktop.Notifications"
	const CORVID_DBUS_OBJECT = "/sh/cxl/Corvid"
	const CORVID_DBUS_NAME = "sh.cxl.Corvid"

	DEFAULT_EXPIRATION = defaultExpiration
	SORT_DIRECTION = sortDirection

	notifications := notificationStack{
		mutex:         &sync.Mutex{},
		notifications: make(map[uint32]notification),
		nextId:        1,
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	err = startDBusServer(
		conn,
		corvidServer{
			notifications: &notifications,
			conn:          conn,
			object:        CORVID_DBUS_OBJECT,
			name:          CORVID_DBUS_NAME,
		},
		CORVID_DBUS_OBJECT,
		CORVID_DBUS_NAME,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = startDBusServer(
		conn,
		notifServer{
			notifications: &notifications,
			conn:          conn,
			object:        NOTIF_DBUS_OBJECT,
			name:          NOTIF_DBUS_NAME,
		},
		NOTIF_DBUS_OBJECT,
		NOTIF_DBUS_NAME,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[]")
}

func startDBusServer(conn *dbus.Conn, v interface{}, object dbus.ObjectPath, name string) error {
	conn.Export(v, object, name)

	reply, err := conn.RequestName(name, dbus.NameFlagReplaceExisting|dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("'%s' already taken", name)
	}

	log.Printf("connected to dbus as %s @ %s", name, object)

	return nil
}
