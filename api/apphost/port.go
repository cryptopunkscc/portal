package apphost

import (
	"slices"
	"strings"
)

type Port []string

func (p Port) Base() string {
	if len(p) == 0 {
		return ""
	}
	return p[0]
}

func (p Port) Name() string {
	if len(p) < 2 {
		return ""
	}
	return strings.Join(p[1:], ".")
}

func NewPort(chunks ...string) Port {
	return Port{}.Add(chunks...)
}

func (p Port) Add(chunks ...string) Port {
	chunks = slices.DeleteFunc(chunks, func(chunk string) bool { return chunk == "" })
	return append(p, chunks...)
}

func (p Port) String() string {
	return strings.Join(p, ".")
}

func (p Port) ParseCmd(query string) (s string, ok bool) {
	s = strings.TrimPrefix(query, p.String())
	s = strings.TrimPrefix(s, ".")
	ok = len(s) < len(query)
	return
}

func FormatPort(chunk string, chunks ...string) string {
	return NewPort(chunk).Add(chunks...).String()
}
