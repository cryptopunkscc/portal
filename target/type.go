package target

import "strconv"

const (
	TypeAny      = Type(0x0)
	TypeBackend  = Type(0x1)
	TypeFrontend = Type(0x2)
	TypeDev      = Type(0x4)
	TypeBundle   = Type(0x8)
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
