package apphost

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/mem"
)

type Cached interface {
	Client
	Cache
	Interrupt()
}

type Cache interface {
	Connections() mem.ReadCache[Conn]
	Events() *sig.Queue[Event]
}

type Event struct {
	Type  EventType
	Query string
	Ref   string
}

type EventType int

const (
	EventConnect EventType = iota
	EventDisconnect
	EventRegister
	EventUnregister
)
