package stream

var End = end{}

type end struct{}

func (e end) Error() string                { return "end" }
func (e end) MarshalJSON() ([]byte, error) { return nil, e }
