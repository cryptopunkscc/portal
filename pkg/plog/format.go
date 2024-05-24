package plog

import (
	"fmt"
	"os"
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
	format string
	time   string
	level  string
}

func NewFormatter(format, time, level string) *Formatter {
	return &Formatter{format, time, level}
}

func (f Formatter) Bytes(l Log) []byte {
	return []byte(f.String(l))
}

func (f Formatter) String(l Log) (line string) {
	line = fmt.Sprintf(f.format,
		l.Time.Format(f.time),
		l.Pid,
		f.level[l.Level],
		l.Scopes,
		l.Message,
	)
	if l.Level <= Fatal {
		line += fmt.Sprintf("%s", l.Stack)
	}
	return
}
