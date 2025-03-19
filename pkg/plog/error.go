package plog

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

func Err(err error, msg ...any) ErrStack {
	return wrapErrStack(err, msg...)
}

func Errorf(format string, args ...any) ErrStack {
	return wrapErrStack(fmt.Errorf(format, args...))
}

func TraceErr(errPtr *error) {
	if errPtr != nil && *errPtr != nil {
		*errPtr = wrapErrStack(*errPtr)
	}
}

type ErrStack struct {
	error
	stack []byte
}

func (e ErrStack) Unwrap() error { return e.error }

func (e ErrStack) Error() string {
	return e.error.Error()
}

func (e ErrStack) Stack() []byte {
	return e.stack
}

func wrapErrStack(err error, msg ...any) (out ErrStack) {
	defer func() {
		if len(msg) == 0 {
			return
		}
		var m string
		if f, ok := (msg[0]).(string); ok && strings.Contains(f, "%") {
			m = fmt.Sprintf(f, msg[1:]...)
		} else {
			m = fmt.Sprint(msg...)
		}
		if out.error != nil {
			out.error = fmt.Errorf("%s: %w", m, err)
		} else if m != "" {
			out.error = errors.New(m)
		}
	}()

	//goland:noinspection GoTypeAssertionOnErrors
	out, ok := err.(ErrStack)
	if ok {
		return out
	}

	stack := debug.Stack()
	chunks := bytes.SplitN(stack, []byte("\n"), 8)
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
