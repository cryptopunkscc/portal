package bind

import (
	"bufio"
	"sync"

	"github.com/cryptopunkscc/astrald/lib/apphost"
)

type cache struct {
	sync.RWMutex
	values map[string]cachedConn
}

func (c *cache) init() {
	if c.values == nil {
		c.values = make(map[string]cachedConn)
	}
}

func (c *cache) interrupt() {
	c.Lock()
	defer c.Unlock()
	for _, value := range c.values {
		_ = value.Close()
	}
	c.values = nil
}

func (c *cache) set(ac apphost.Conn) conn {
	c.Lock()
	defer c.Unlock()
	c.init()
	cc := conn{ac, *bufio.NewReader(ac)}
	c.values[ac.Query().Nonce.String()] = cachedConn{cc, c}
	return cc
}

func (c *cache) delete(id string) {
	c.Lock()
	defer c.Unlock()
	delete(c.values, id)
}

func (c *cache) get(id string) (conn, bool) {
	c.RLock()
	defer c.RUnlock()
	cc, ok := c.values[id]
	return cc.conn, ok
}

type cachedConn struct {
	conn
	*cache
}

func (c *cachedConn) Read(p []byte) (n int, err error) {
	if n, err = c.Reader.Read(p); err != nil {
		_ = c.Close()
	}
	return
}

func (c *cachedConn) Write(p []byte) (n int, err error) {
	if n, err = c.Conn.Write(p); err != nil {
		_ = c.Close()
	}
	return
}

func (c *cachedConn) ReadString(delim byte) (s string, err error) {
	if s, err = c.Reader.ReadString(delim); err != nil {
		_ = c.Close()
	}
	return
}

func (c *cachedConn) Close() error {
	c.cache.delete(c.Conn.Query().Nonce.String())
	return c.Conn.Close()
}
