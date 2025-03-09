package plog

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
)

type ErrStack struct {
	error
	stack []byte
}

func Err(err error) (out ErrStack) {
	if errors.As(err, &out) {
		return
	}
	stack := debug.Stack()
	chunks := bytes.SplitN(stack, []byte("\n"), 6)
	if len(chunks) > 0 {
		stack = chunks[0]
	}
	if len(chunks) > 1 {
		stack = append(stack, '\n')
		stack = append(stack, chunks[len(chunks)-1]...)
	}
	return ErrStack{
		error: err,
		stack: stack,
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
