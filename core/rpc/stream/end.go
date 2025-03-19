package stream

var End = end{}

type end struct{}

func (e end) Error() string { return "end" }
