package srv

import (
	"github.com/godbus/dbus/v5"
)

type corvidServer server

func (s corvidServer) Dismiss(id uint32) (e *dbus.Error) {
	// log.Print("Dismiss called")
	server(s).close(id, CloseReasonDismissed)
	server(s).output()
	return nil
}

func (s corvidServer) DismissAll() (e *dbus.Error) {
	// log.Print("DismissAll called")
	for _, notification := range s.notifications.notifications {
		server(s).close(notification.Id, CloseReasonDismissed)
	}

	server(s).output()
	return nil
}
