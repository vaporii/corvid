package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/godbus/dbus/v5"
)

const DEFAULT_EXPIRATION = 5000
const SORT_DIRECTION = 1 // 1 = newest first, -1 = oldest first
const DBUS_OBJECT = "/org/freedesktop/Notifications"
const DBUS_NAME = "org.freedesktop.Notifications"

var conn *dbus.Conn

type notificationStack = struct {
	mutex         *sync.Mutex
	notifications map[uint32]notification
	nextId        uint32
}

var notifications = notificationStack{
	mutex:         &sync.Mutex{},
	notifications: make(map[uint32]notification),
	nextId:        1,
}

func output() {
	arr := make([]notification, len(notifications.notifications))

	i := 0
	for _, notification := range notifications.notifications {
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

func main() {
	var err error
	conn, err = dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	conn.Export(server{}, DBUS_OBJECT, DBUS_NAME)

	reply, err := conn.RequestName(DBUS_NAME, dbus.NameFlagReplaceExisting|dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatal(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalf("'%s' already taken", DBUS_NAME)
	}

	log.Print("connected to dbus")
	select {}
}
