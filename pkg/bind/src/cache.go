package bind

import (
	"bufio"
	"sync"

	"github.com/cryptopunkscc/astrald/lib/apphost"
)

type cache struct {
	sync.RWMutex
	values map[string]conn
}

func (c *cache) init() {
	if c.values == nil {
		c.values = make(map[string]conn)
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
	cc := conn{ac, *bufio.NewReader(ac), c}
	c.values[ac.Query().Nonce.String()] = cc
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
	return cc, ok
}

type conn struct {
	apphost.Conn
	bufio.Reader
	*cache
}

func (c *conn) Read(p []byte) (n int, err error) {
	if n, err = c.Reader.Read(p); err != nil {
		_ = c.Close()
	}
	return
}

func (c *conn) Write(p []byte) (n int, err error) {
	if n, err = c.Conn.Write(p); err != nil {
		_ = c.Close()
	}
	return
}

func (c *conn) ReadString(delim byte) (s string, err error) {
	if s, err = c.Reader.ReadString(delim); err != nil {
		_ = c.Close()
	}
	return
}

func (c *conn) Close() error {
	c.cache.delete(c.Conn.Query().Nonce.String())
	return c.Conn.Close()
}
