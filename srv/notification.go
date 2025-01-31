package srv

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

type hint struct {
	dbus.Variant
}

func (h hint) MarshalJSON() ([]byte, error) {
	//TODO: find a better way lol
	switch h.Signature().String()[0] {
	case 'y':
		return json.Marshal(h.Value().(uint8))
	case 'b':
		return json.Marshal(h.Value().(bool))
	case 'n':
		return json.Marshal(h.Value().(int16))
	case 'q':
		return json.Marshal(h.Value().(uint16))
	case 'i':
		return json.Marshal(h.Value().(int32))
	case 'u':
		return json.Marshal(h.Value().(uint32))
	case 'x':
		return json.Marshal(h.Value().(int64))
	case 't':
		return json.Marshal(h.Value().(uint64))
	case 'd':
		return json.Marshal(h.Value().(float64))
	case 's':
		return json.Marshal(h.Value().(string))
	default:
		panic("Impossible type")
	}
}

type closeReason uint32

const (
	CloseReasonExpire    closeReason = 1
	CloseReasonDismissed             = iota
	CloseReasonClosed                = iota
	CloseReasonOther                 = iota
)

type notification struct {
	Id         uint32            `json:"id"`
	AppName    string            `json:"app_name"`
	AppIcon    string            `json:"app_icon"`
	Summary    string            `json:"summary"`
	Body       string            `json:"body"`
	Actions    map[string]string `json:"actions"`
	Hints      map[string]hint   `json:"hints"`
	Timestamp  int64             `json:"timestamp"`
	Expiration int32             `json:"expiration"`
	Image      string            `json:"image"`
	timer      *time.Timer
}

type notificationStack = struct {
	mutex         *sync.Mutex
	notifications map[uint32]notification
	nextId        uint32
}
