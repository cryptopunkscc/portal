package rpc

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

type TestGoService struct {
	*rpc.App
}

func NewTestGoService(port string) *TestGoService {
	app := rpc.NewApp(port)
	app.Logger(plog.New().Type(app))

	app.RouteFunc("request.func1", Function1)
	app.RouteFunc("request.func2", Function2)
	app.RouteFunc("request.func3", Function3)
	app.RouteFunc("request.func4", Function4)

	app.Routes("flow*")
	app.Func("flow.func1", Function1)
	app.Func("flow.func2", Function2)
	app.Func("flow.func3", Function3)
	app.Func("flow.func4", Function4)

	return &TestGoService{App: app}
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

func Function3(struct1 *TestStruct2) *TestStruct2 {
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
	B bool    `json:"b"`
	I int     `json:"i"`
	F float64 `json:"f"`
	S string  `json:"s"`
}

type TestStruct2 struct {
	TestStruct1 `json:"struct1"`
}
