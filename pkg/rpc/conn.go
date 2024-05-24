package rpc

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"io"
)

type Conn interface {
	io.WriteCloser
	ByteScannerReader
	Logger(logger plog.Logger)
	Copy() Conn
	Call(method string, value any) (err error)
	Encode(value any) (err error)
	Decode(value any) (err error)
	Flush()
}

func Call(conn Conn, name string, args ...any) (err error) {
	var payload any
	if len(args) > 0 {
		payload = args
	}
	return conn.Call(name, payload)
}

func Decode[R any](conn Conn) (r R, err error) {
	err = conn.Decode(&r)
	return
}

func Await(conn Conn) (err error) {
	var r interface{}
	err = conn.Decode(&r)
	return
}

func Command(conn Conn, method string, args ...any) (err error) {
	conn = conn.Copy()
	defer conn.Flush()
	if err = Call(conn, method, args...); err != nil {
		return
	}
	if err = Await(conn); errors.Is(err, io.EOF) {
		err = nil
	}
	return
}

func Query[R any](conn Conn, method string, args ...any) (r R, err error) {
	conn = conn.Copy()
	defer conn.Flush()
	if err = Call(conn, method, args...); err == nil {
		r, err = Decode[R](conn)
	}
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	return
}

func Subscribe[R any](conn Conn, method string, args ...any) (c <-chan R, err error) {
	conn = conn.Copy()
	if err = Call(conn, method, args...); err != nil {
		return
	}
	cc := make(chan R)
	go func() {
		defer close(cc)
		defer conn.Flush()
		var r R
		for {
			if err = conn.Decode(&r); err != nil {
				return
			}
			cc <- r
		}
	}()
	c = cc
	return
}
