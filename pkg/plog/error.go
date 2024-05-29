package plog

import (
	"fmt"
	"runtime/debug"
)

type ErrStack struct {
	error
	stack []byte
}

func Err(err error) ErrStack {
	return ErrStack{
		error: err,
		stack: debug.Stack(),
	}
}

func Errorf(format string, args ...any) ErrStack {
	return ErrStack{
		error: fmt.Errorf(format, args...),
		stack: debug.Stack(),
	}
}

func (e ErrStack) Msgf(format string, args ...any) ErrStack {
	msg := fmt.Sprintf(format, args...)
	e.error = fmt.Errorf("%s: %v", msg, e.error)
	return e
}

func (e ErrStack) Error() string {
	return e.error.Error()
}

func (e ErrStack) Stack() []byte {
	return e.stack
}
