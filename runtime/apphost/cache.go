package apphost

import (
	"bufio"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/sig"
	. "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/google/uuid"
)

type cache struct {
	connections mem.Cache[Conn]
	listeners   mem.Cache[Listener]
	events      *sig.Queue[Event]
}

func newCache() *cache {
	return &cache{
		connections: mem.NewCache[Conn](),
		listeners:   mem.NewCache[Listener](),
		events:      &sig.Queue[Event]{},
	}
}

func (c *cache) Connections() mem.ReadCache[Conn]   { return c.connections }
func (c *cache) Listeners() mem.ReadCache[Listener] { return c.listeners }
func (c *cache) Events() *sig.Queue[Event]          { return c.events }

func (c *cache) inConn(ac *astral.Conn, err error) (Conn, error)  { return c.setConn(ac, err, true) }
func (c *cache) outConn(ac *astral.Conn, err error) (Conn, error) { return c.setConn(ac, err, false) }
func (c *cache) setConn(ac *astral.Conn, err error, in bool) (Conn, error) {
	if err != nil {
		return nil, err
	}
	out := &conn{
		Conn:  ac,
		cache: c,
		buf:   bufio.NewReader(ac),
		ref:   uuid.New().String(),
		in:    in,
	}
	c.connections.Set(out.ref, out)
	c.events.Push(Event{Type: Connect, Port: out.Query(), Ref: out.ref})
	return out, nil
}
func (c *cache) deleteConn(conn Conn) {
	c.connections.Delete(conn.Ref())
	c.events.Push(Event{Type: Disconnect, Port: conn.Query(), Ref: conn.Ref()})
}
func (c *cache) setListener(port string, al *astral.Listener, err error) (Listener, error) {
	if err != nil {
		return nil, err
	}
	ll := &listener{al, c, port}
	c.listeners.Set(ll.port, ll)
	c.events.Push(Event{Type: Register, Port: port})
	return ll, nil
}
func (c *cache) deleteListener(port string) {
	c.listeners.Delete(port)
	c.events.Push(Event{Type: Unregister, Port: port})
	return
}
