package target

import "strconv"

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
		if i, err := strconv.Atoi(args[0]); err == nil {
			typ = Type(i)
			return
		}
	}
	return def
}
