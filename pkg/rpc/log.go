package rpc

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"io"
	"strings"
)

type ConnLogger struct {
	io.ReadWriteCloser
	plog.Logger
}

func NewConnLogger(conn io.ReadWriteCloser, logger plog.Logger) *ConnLogger {
	return &ConnLogger{
		ReadWriteCloser: conn,
		Logger:          logger,
	}
}

func (cl *ConnLogger) Read(b []byte) (n int, err error) {
	n, err = cl.ReadWriteCloser.Read(b)
	if n > 0 {
		cl.Println("<", strings.Trim(string(b[:n]), "\n"))
	}
	return
}

func (cl *ConnLogger) Write(b []byte) (n int, err error) {
	n, err = cl.ReadWriteCloser.Write(b)
	if n > 0 {
		cl.Println(">", strings.Trim(string(b[:n]), "\n"))
	}
	return
}
