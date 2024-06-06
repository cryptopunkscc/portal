package rpc

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

type TestGoService struct {
	*rpc.App
}

func NewTestGoService(port string) *TestGoService {
	app := rpc.NewApp(port)
	app.RouteFunc("func1", Function1)
	app.RouteFunc("func2", Function2)
	app.RouteFunc("func3", Function3)
	app.RouteFunc("func4", Function4)
	app.Routes("flow.*")
	app.Func("flow.func1", Function1)
	app.Func("flow.func2", Function2)
	return &TestGoService{App: app}
}

func Function1(arg string, err error) (string, error) {
	return arg, err
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
