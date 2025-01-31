package srv

import (
	"log"

	"github.com/godbus/dbus/v5"
)

type corvidServer server

func (s corvidServer) Test(param uint32) (e *dbus.Error) {
	log.Printf("Test called %d", param)
	return nil
}

// func (s corvidServer) GetCapabilities() (e *dbus.Error) {
// 	log.Print("GetCapabilities called")
// 	return nil
// }

// func (s corvidServer) GetServerInformation() (name, vendor, version, specVersion string, e *dbus.Error) {
// 	// log.Print("GetServerInformation called")
// 	return "corvid", "CartConnoisseur", "0.1.0", "1.2", nil
// }

// func (s corvidServer) CloseNotification(id uint32) (e *dbus.Error) {
// 	// log.Printf("CloseNotification called: %d", id)
// 	notification, ok := notifications.notifications[id]
// 	if ok {
// 		notification.close(CloseReasonClosed)
// 	}

// 	return nil
// }

// func (s corvidServer) Notify(appName string, replacesId uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (id uint32, e *dbus.Error) {
// 	// log.Print("Notify called")

// }
