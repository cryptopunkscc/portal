package plog

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"time"
)

type logger struct {
	out Output
	Log
}

func New() Logger {
	return &logger{
		out: DefaultOutput,
		Log: Log{Pid: os.Getpid()},
	}
}

const key = "plog"

func Get(ctx context.Context) Logger {
	return ctx.Value(key).(Logger)
}

func (l logger) Set(ctx *context.Context) Logger {
	*ctx = context.WithValue(*ctx, key, l)
	return l
}

func (l logger) Out(output Output) Logger {
	l.out = output
	return l
}

func (l logger) Scope(format string, args ...any) Logger {
	if args != nil {
		l.Scopes = append(l.Scopes, fmt.Sprintf(format, args...))
		return l
	}
	l.Scopes = append(l.Scopes, format)
	return l
}

func (l logger) Type(a any) Logger {
	l.Scopes = append(l.Scopes, reflect.TypeOf(a))
	return l
}

func (l logger) Any(a any) Logger {
	l.Scopes = append(l.Scopes, a)
	return l
}

func (l logger) P() Logger {
	l.Level = Panic
	return l
}

func (l logger) F() Logger {
	l.Level = Fatal
	return l
}

func (l logger) E() Logger {
	l.Level = Error
	return l
}

func (l logger) W() Logger {
	l.Level = Warning
	return l
}

func (l logger) I() Logger {
	l.Level = Info
	return l
}

func (l logger) D() Logger {
	l.Level = Debug
	return l
}

func (l logger) Msg(message string) Logger {
	l.Message = message
	return l
}

func (l logger) Printf(format string, a ...any) {
	l.Message = fmt.Sprintf(format, a...)
	l.Flush()
}

func (l logger) Println(a ...any) {
	l.Message = fmt.Sprintln(a...)
	l.Flush()
}

func (l logger) Flush() {
	l.Time = time.Now()
	if l.Level <= Fatal {
		l.Stack = debug.Stack()
	}
	l.out(l.Log)
	if l.Level == Panic {
		os.Exit(1)
	}
}
