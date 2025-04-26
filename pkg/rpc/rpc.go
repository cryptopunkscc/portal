package rpc

type Rpc interface {
	Conn(target, query string) (Conn, error)
	Request(target string, query ...string) Conn
}
