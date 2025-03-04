package apphost

import (
	"github.com/cryptopunkscc/portal/api/apphost"
)

func Cached(client apphost.Client) apphost.Cached {
	return &cached{Client: client, cache: newCache()}
}

type cached struct {
	apphost.Client
	*cache
}

func (a cached) Interrupt() {
	for _, closer := range a.listeners.Release() {
		_ = closer.Close()
	}
	for _, closer := range a.connections.Release() {
		_ = closer.Close()
	}
}

func (a cached) Query(target string, method string, args any) (conn apphost.Conn, err error) {
	return a.setConn(a.Client.Query(target, method, args))
}

func (a cached) Register() (l apphost.Listener, err error) {
	ll, err := a.Client.Register()
	return a.setListener(ll, err)
}

type cachedListener struct {
	apphost.Listener
	cache *cache
}

func (l *cachedListener) Next() (q apphost.PendingQuery, err error) {
	qq := cachedQuery{cache: l.cache}
	if qq.PendingQuery, err = l.Listener.Next(); err != nil {
		return
	}
	return &qq, nil
}

func (l *cachedListener) Close() (err error) {
	err = l.Listener.Close()
	l.cache.deleteListener(l.String())
	return
}

type cachedQuery struct {
	apphost.PendingQuery
	*cache
}

func (q *cachedQuery) Accept() (apphost.Conn, error) {
	return q.setConn(q.PendingQuery.Accept())
}

type cachedConn struct {
	apphost.Conn
	*cache
}

func (c *cachedConn) Read(p []byte) (n int, err error) {
	if n, err = c.Conn.Read(p); err != nil {
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
	if s, err = c.Conn.ReadString(delim); err != nil {
		_ = c.Close()
	}
	return
}

func (c *cachedConn) Close() error {
	c.cache.deleteConn(c)
	return c.Conn.Close()
}
