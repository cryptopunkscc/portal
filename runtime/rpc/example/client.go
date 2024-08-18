package main

import (
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

type apiClient struct {
	conn rpc.Conn
}

func NewApiClient(conn rpc.Conn) Api {
	return &apiClient{conn}
}

func (a apiClient) Method(b bool, i int, s string) {
	_ = rpc.Command(a.conn, "method", b, i, s)
}

func (a apiClient) Method1(b bool) (err error) {
	return rpc.Command(a.conn, "method1", b)
}

func (a apiClient) Method2(arg *Arg) (Arg, error) {
	return rpc.Query[Arg](a.conn, "method2", arg)
}

func (a apiClient) Method2S() (string, error) {
	return rpc.Query[string](a.conn, "method2S")
}

func (a apiClient) Method2B() (bool, error) {
	return rpc.Query[bool](a.conn, "method2B")
}

func (a apiClient) MethodC() (<-chan Arg, error) {
	return rpc.Subscribe[Arg](a.conn, "methodC")
}
