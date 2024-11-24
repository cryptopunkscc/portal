package port

import (
	"slices"
	"strings"
)

type Port []string

func New(chunks ...string) Port {
	return Prefix().Add(chunks...)
}

func (p Port) Add(chunks ...string) Port {
	chunks = slices.DeleteFunc(chunks, func(chunk string) bool { return chunk == "" })
	return append(p, chunks...)
}

func (p Port) String() string {
	return strings.Join(p, ".")
}

func (p Port) ParseCmd(query string) (s string) {
	s = strings.TrimPrefix(query, p.String())
	s = strings.TrimPrefix(s, ".")
	return
}
