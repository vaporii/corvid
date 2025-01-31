package srv

import (
	"fmt"
	"log"
	"sync"

	"github.com/godbus/dbus/v5"
)

const DEFAULT_EXPIRATION = 5000
const SORT_DIRECTION = 1 // 1 = newest first, -1 = oldest first

type server = struct {
	conn   *dbus.Conn
	object dbus.ObjectPath
	name   string
}

func Start() {
	const NOTIF_DBUS_OBJECT = "/org/freedesktop/Notifications"
	const NOTIF_DBUS_NAME = "org.freedesktop.Notifications"
	const CORVID_DBUS_OBJECT = "/sh/cxl/Corvid"
	const CORVID_DBUS_NAME = "sh.cxl.Corvid"

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
			conn:   conn,
			object: CORVID_DBUS_OBJECT,
			name:   CORVID_DBUS_NAME,
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
			server: server{
				conn:   conn,
				object: NOTIF_DBUS_OBJECT,
				name:   NOTIF_DBUS_NAME,
			},
		},
		NOTIF_DBUS_OBJECT,
		NOTIF_DBUS_NAME,
	)
	if err != nil {
		log.Fatal(err)
	}
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
