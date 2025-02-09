package plog

import (
	"fmt"
	"runtime/debug"
)

var ErrorStackTrace = false

type ErrStack struct {
	error
	stack []byte
}

func Err(err error) error {
	return ErrStack{
		error: err,
		stack: debug.Stack(),
	}
}

func Errorf(format string, args ...any) (err error) {
	return Err(fmt.Errorf(format, args...))
}

func (e ErrStack) Error() string {
	return e.error.Error()
}

func (e ErrStack) Stack() []byte {
	return e.stack
}
