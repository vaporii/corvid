package main

import (
	"log"

	"github.com/godbus/dbus/v5"
)

func main() {
	const CORVID_DBUS_OBJECT = "/sh/cxl/Corvid"
	const CORVID_DBUS_NAME = "sh.cxl.Corvid"

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	call := conn.Object(CORVID_DBUS_NAME, CORVID_DBUS_OBJECT).Call(CORVID_DBUS_NAME+".Test", 0, uint32(13))
	if call.Err != nil {
		log.Fatal(call.Err)
	}
}
