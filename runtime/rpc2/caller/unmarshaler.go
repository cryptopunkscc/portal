package caller

type Unmarshaler interface {
	Unmarshal(data []byte, args []any) error
}

type Unmarshal func(data []byte, args []any) error

func (u Unmarshal) Unmarshal(data []byte, args []any) error { return u(data, args) }
