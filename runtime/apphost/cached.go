package apphost

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
)

type cached struct {
	apphost.Client
	*cache
}

func Cached(client apphost.Client) apphost.Cached {
	return &cached{Client: client, cache: newCache()}
}

func (a cached) Interrupt() {
	for _, closer := range a.listeners.Release() {
		_ = closer.Close()
	}
	for _, closer := range a.connections.Release() {
		_ = closer.Close()
	}
}

func (a cached) Query(remoteID id.Identity, query string) (conn apphost.Conn, err error) {
	return a.setConn(a.Client.Query(remoteID, query))
}
func (a cached) Register(service string) (l apphost.Listener, err error) {
	ll, err := a.Client.Register(service)
	return a.setListener(service, ll, err)
}

type cachedListener struct {
	apphost.Listener
	cache *cache
}

func (l *cachedListener) Next() (q apphost.QueryData, err error) {
	qq := cachedQuery{cache: l.cache}
	if qq.QueryData, err = l.Listener.Next(); err != nil {
		return
	}
	return &qq, nil
}
func (l *cachedListener) QueryCh() (c <-chan apphost.QueryData) {
	out := make(chan apphost.QueryData)
	go func() {
		defer close(out)
		for data := range l.Listener.QueryCh() {
			out <- &cachedQuery{data, l.cache}
		}
	}()
	return out
}
func (l *cachedListener) Close() (err error) {
	err = l.Listener.Close()
	l.cache.deleteListener(l.Port())
	return
}

type cachedQuery struct {
	apphost.QueryData
	*cache
}

func (q *cachedQuery) Accept() (apphost.Conn, error) {
	return q.setConn(q.QueryData.Accept())
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
