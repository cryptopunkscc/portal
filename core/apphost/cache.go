package apphost

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
)

func newCache() *cache {
	return &cache{
		connections: mem.NewCache[apphost.Conn](),
		listeners:   mem.NewCache[apphost.Listener](),
		events:      &sig.Queue[apphost.Event]{},
	}
}

type cache struct {
	connections mem.Cache[apphost.Conn]
	listeners   mem.Cache[apphost.Listener]
	events      *sig.Queue[apphost.Event]
}

func (c *cache) Connections() mem.ReadCache[apphost.Conn]   { return c.connections }
func (c *cache) Listeners() mem.ReadCache[apphost.Listener] { return c.listeners }
func (c *cache) Events() *sig.Queue[apphost.Event]          { return c.events }

func (c *cache) setConn(ac apphost.Conn, err error) (apphost.Conn, error) {
	if err != nil {
		return nil, err
	}
	c.connections.Set(ac.Ref(), ac)
	c.events.Push(apphost.Event{Type: apphost.EventConnect, Query: ac.Query(), Ref: ac.Ref()})
	return ac, nil
}

func (c *cache) deleteConn(conn apphost.Conn) {
	c.connections.Delete(conn.Ref())
	c.events.Push(apphost.Event{Type: apphost.EventDisconnect, Query: conn.Query(), Ref: conn.Ref()})
}

func (c *cache) setListener(al apphost.Listener, err error) (apphost.Listener, error) {
	if err != nil {
		return nil, err
	}
	ll := &cachedListener{al, c}
	ll.String()
	c.listeners.Set(ll.String(), ll)
	c.events.Push(apphost.Event{Type: apphost.EventRegister})
	return ll, nil
}

func (c *cache) deleteListener(port string) {
	c.listeners.Delete(port)
	c.events.Push(apphost.Event{Type: apphost.EventUnregister})
	return
}
