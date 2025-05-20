package plog

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"
)

var Verbosity = Panic

type logger struct {
	out    Output
	module string
	Log
}

func (l logger) Copy() Logger {
	ll := l
	ll.Scopes = slices.Clone(l.Scopes)
	return ll
}

var Default = New()

func P() Logger                               { return Default.P() }
func D() Logger                               { return Default.D() }
func Type(a any) Logger                       { return Default.Type(a) }
func Scope(format string, args ...any) Logger { return Default.Scope(format, args...) }
func Println(a ...any)                        { Default.Println(a...) }
func Printf(format string, args ...any)       { Default.Printf(format, args...) }

func New() Logger {
	return logger{
		out:    DefaultOutput,
		module: api.Module,
		Log: Log{
			Pid:   os.Getpid(),
			Level: Info,
		},
	}
}

const key = "plog"

func Get(ctx context.Context) Logger {
	if ctx != nil {
		if l, ok := ctx.Value(key).(Logger); ok {
			return l.Copy()
		}
	}
	return Default
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
		l.Scopes = append(l.Scopes[:], fmt.Sprintf(format, args...))
		return l
	}
	l.Scopes = append(l.Scopes[:], format)
	return l
}

func (l logger) Type(a any) Logger {
	t := reflect.TypeOf(a)
	s := strings.Replace(t.String(), l.module, "", -1)
	l.Scopes = append(l.Scopes[:], s)
	return l
}

func (l logger) Any(a any) Logger {
	l.Scopes = append(l.Scopes[:], fmt.Sprint(a))
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
	if l.Level > Verbosity {
		return l
	}
	l.Message += message
	return l
}

func (l logger) Printf(format string, args ...any) {
	if l.Level > Verbosity {
		return
	}
	l.Message += fmt.Sprintf(format, args...) + "\n"
	l.appendErrors(args...).Flush()
}
func (l logger) Println(a ...any) {
	if l.Level > Verbosity {
		return
	}
	l.Message += fmt.Sprintln(a...)
	l.appendErrors(a...).Flush()
}

func (l logger) Flush() {
	if l.Level > Verbosity {
		return
	}
	l.Time = time.Now()
	l.out(l.Log)
	if l.Level == Panic {
		os.Exit(1)
	}
}

func (l logger) appendErrors(args ...any) logger {
	for _, arg := range args {
		switch e := arg.(type) {
		case []error:
			for _, err := range e {
				l.Errors = append(l.Errors, err)
			}
		case error:
			l.Errors = append(l.Errors, e)
		}
	}
	return l
}
