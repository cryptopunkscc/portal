package rpc

import (
	"errors"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

type TestGoService struct {
	*apphost.Router
}

func NewTestGoService(p string) *TestGoService {
	handlers := cmd.Handlers{
		{Func: Function1, Name: "func1"},
		{Func: Function2, Name: "func2"},
		{Func: Function3, Name: "func3"},
		{Func: Function4, Name: "func4"},
	}
	root := cmd.Handler{
		Name: p,
		Sub: cmd.Handlers{
			{
				Name: "request",
				Sub:  handlers,
			},
			{
				Name: "flow",
				Sub:  handlers,
				Func: apphost.RouteAll,
			},
		},
	}
	return &TestGoService{
		Router: apphost.NewRouter(root),
	}
}

func Function1(msg string, fail bool) (s string, err error) {
	s = msg
	if fail {
		err = errors.New(msg)
	}
	return
}

func Function2(
	b bool,
	i int,
	f float64,
	s string,
) (
	bool,
	int,
	float64,
	string,
) {
	return b, i, f, s
}

func Function3(struct1 *TestStruct1) *TestStruct1 {
	return struct1
}

func Function4(
	b bool,
	i int,
	f float64,
	s string,
) TestStruct1 {
	return TestStruct1{b, i, f, s}
}

type TestStruct1 struct {
	B bool    `query:"b" json:"b"`
	I int     `query:"i" json:"i"`
	F float64 `query:"f" json:"f"`
	S string  `query:"s" json:"s"`
}
