package caller

type Unmarshaler interface {
	Unmarshal(data []byte, args []any) error
	Score(data []byte) (score uint)
}

type Unmarshal func(data []byte, args []any) error

func (u Unmarshal) Unmarshal(data []byte, args []any) error { return u(data, args) }
func (u Unmarshal) Score([]byte) uint                       { return 1 }
