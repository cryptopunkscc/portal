package port

import (
	"strings"
)

type Port []string

func New(chunks ...string) Port {
	return Prefix().Add(chunks...)
}

func (p Port) Add(chunks ...string) Port {
	return append(p, chunks...)
}

func (p Port) String() string {
	return strings.Join(p, ".")
}
