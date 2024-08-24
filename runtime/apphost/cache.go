package apphost

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
)

type cache struct {
	connections mem.Cache[apphost.Conn]
	listeners   mem.Cache[apphost.Listener]
	events      *sig.Queue[apphost.Event]
}

func newCache() *cache {
	return &cache{
		connections: mem.NewCache[apphost.Conn](),
		listeners:   mem.NewCache[apphost.Listener](),
		events:      &sig.Queue[apphost.Event]{},
	}
}

func (c *cache) Connections() mem.ReadCache[apphost.Conn]   { return c.connections }
func (c *cache) Listeners() mem.ReadCache[apphost.Listener] { return c.listeners }
func (c *cache) Events() *sig.Queue[apphost.Event]          { return c.events }

func (c *cache) setConn(ac apphost.Conn, err error) (apphost.Conn, error) {
	if err != nil {
		return nil, err
	}
	ac = &cachedConn{
		Conn:  ac,
		cache: c,
	}
	c.connections.Set(ac.Ref(), ac)
	c.events.Push(apphost.Event{Type: apphost.Connect, Port: ac.Query(), Ref: ac.Ref()})
	return ac, nil
}

func (c *cache) deleteConn(conn apphost.Conn) {
	c.connections.Delete(conn.Ref())
	c.events.Push(apphost.Event{Type: apphost.Disconnect, Port: conn.Query(), Ref: conn.Ref()})
}

func (c *cache) setListener(port string, al apphost.Listener, err error) (apphost.Listener, error) {
	if err != nil {
		return nil, err
	}
	ll := &cachedListener{al, c}
	c.listeners.Set(ll.Port(), ll)
	c.events.Push(apphost.Event{Type: apphost.Register, Port: port})
	return ll, nil
}

func (c *cache) deleteListener(port string) {
	c.listeners.Delete(port)
	c.events.Push(apphost.Event{Type: apphost.Unregister, Port: port})
	return
}
