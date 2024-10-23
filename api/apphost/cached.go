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
	Listeners() mem.ReadCache[Listener]
	Events() *sig.Queue[Event]
}

type Event struct {
	Type EventType
	Port string
	Ref  string
}

type EventType int

const (
	Connect EventType = iota
	Disconnect
	Register
	Unregister
)
