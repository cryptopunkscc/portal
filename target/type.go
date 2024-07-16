package target

import "fmt"

const (
	TypeAny     Type = 0
	TypeBackend Type = 1 << (iota - 1)
	TypeFrontend
	TypeDev
	TypeBundle
)

type Type int

func (t Type) Is(p Type) bool {
	return t&p == p
}

func ParseType(def Type, args ...string) (typ Type) {
	if len(args) > 0 {
		if _, err := fmt.Sscan(args[0], &typ); err != nil {
			return
		}
	}
	return def
}
