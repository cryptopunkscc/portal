package main

type Api interface {
	Method(b bool, i int, s string)
	Method1(bool) error
	Method2(arg *Arg) (Arg, error)
	Method2S() (string, error)
	Method2B() (bool, error)
	MethodC() (c <-chan Arg, err error)
}

type Arg struct {
	S   string
	I   int
	Arg *Arg `json:"arg,omitempty"`
}
