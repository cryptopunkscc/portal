package plog

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const DefaultTimeFormat = "2006-01-02 15:04:05.000000"
const DefaultLogFormat = "%s (%d) %c %v: %s"
const DefaultLevels = "PFEWID"

var DefaultFormatter = NewFormatter(DefaultLogFormat, DefaultTimeFormat, DefaultLevels)

var DefaultOutput Output = func(log Log) {
	bytes := DefaultFormatter.Bytes(log)
	_, _ = os.Stderr.Write(bytes)
}

type Formatter struct {
	format    string
	time      string
	level     string
	scopeSize int
}

func NewFormatter(format, time, level string) *Formatter {
	return &Formatter{format: format, time: time, level: level}
}

func (f *Formatter) Bytes(l Log) []byte {
	return []byte(f.String(l))
}

func (f *Formatter) String(l Log) (line string) {
	scopes := strings.Join(l.Scopes, ">")
	//scopesSize := len(scopes)
	//if scopesSize > f.scopeSize {
	//	f.scopeSize = scopesSize
	//} else {
	//	scopes += strings.Repeat(" ", f.scopeSize-scopesSize)
	//}

	line = fmt.Sprintf(f.format,
		l.Time.Format(f.time),
		l.Pid,
		f.level[l.Level],
		"|"+scopes+"|",
		l.Message,
	)
	if len(l.Stack) > 0 {
		line += fmt.Sprintf("%s", l.Stack)
	}
	for _, err := range l.Errors {
		var e ErrStack
		if errors.As(err, &e) {
			line += fmt.Sprintf("\n%s", e.stack)
		}
	}
	return
}
