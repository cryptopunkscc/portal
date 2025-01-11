package rpc

import (
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

var Close = errors.New("close")

type Conn interface {
	io.ReadWriteCloser
	Logger(logger plog.Logger)
	Copy() Conn
	Call(method string, value any) (err error)
	Bytes() ([]byte, error)
	Encode(value any) (err error)
	Decode(value any) (err error)
	Flush()
}

func Call[T any](conn Conn, name string, args ...T) (err error) {
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

func CommandT[T any](conn Conn, method string, args ...T) (err error) {
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
