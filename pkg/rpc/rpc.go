package rpc

type Rpc interface {
	Format(name string) Rpc
	Conn(target, query string) (Conn, error)
	Request(target string, query ...string) Conn
}
