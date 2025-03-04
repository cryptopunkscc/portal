package rpc

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"strings"
)

type Logger struct {
	io.ReadWriteCloser
	plog.Logger
}

func NewConnLogger(conn io.ReadWriteCloser, logger plog.Logger) *Logger {
	return &Logger{
		ReadWriteCloser: conn,
		Logger:          logger,
	}
}

func (cl *Logger) Read(b []byte) (n int, err error) {
	n, err = cl.ReadWriteCloser.Read(b)
	if n > 0 {
		cl.Println("<", strings.Trim(string(b[:n]), "\n"))
	}
	return
}

func (cl *Logger) Write(b []byte) (n int, err error) {
	n, err = cl.ReadWriteCloser.Write(b)
	if n > 0 {
		cl.Println(">", strings.Trim(string(b[:n]), "\n"))
	}
	return
}
