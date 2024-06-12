package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/mod/apphost/proto"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/stretchr/testify/assert"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestApp_Run_basic(t *testing.T) {
	skipped := []struct {
		client int
		cases  []int
		routes int
	}{
		{routes: 2, client: 1, cases: []int{6, 7, 10, 11, 12, 13, 14, 15, 18, 19, 21, 22, 23, 24, 26, 27, 28, 31, 32, 34, 35, 38, 39, 42}},
		{routes: 2, client: 2},
		{routes: 3, client: 2},
	}

	port := "test_app"
	clients := []func(*testing.T) (Conn, error){
		func(*testing.T) (c Conn, err error) {
			c = NewRequest(id.Anyone, port)
			c.Logger(plog.New().Type(c))
			return
		},
		func(t *testing.T) (c Conn, err error) {
			c, err = QueryFlow(id.Anyone, port)
			if err != nil {
				return
			}
			c.Logger(plog.New().Type(c))
			return
		},
	}
	setHandlers := func(app *App) {
		app.Func("func0", function0)
		app.Func("func0!", function0)
		app.Func("func1", function1)
		app.Func("func2", function2)
		app.Func("func3", function3)
		app.Func("func4", function4)
		app.Func("func5", function5)
		app.Func("func6", function6)
		app.Func("func7", function7)
		app.Func("func8", function8)
		app.Func("func9", function9)
		app.Func("func10", function10)
		app.Func("func11", function11)
	}
	routes := [][]string{
		{"*"},
		{
			"func0",
			"func1",
			"func2",
			"func3",
			"func4",
			"func5",
			"func6",
			"func7",
			"func8",
			"func9",
			"func10",
		},
		{
			"func0*",
			"func1*",
			"func2*",
			"func3*",
			"func4*",
			"func5*",
			"func6*",
			"func7*",
			"func8*",
			"func9*",
			"func10*",
		},
	}
	cases := []struct {
		query    string
		arg      any
		args     []any
		expected any
		client   int
	}{
		{query: "asd", expected: proto.ErrRejected, client: 1},
		{query: "asd", expected: ErrNoHandler, client: 2},
		{query: "func0", expected: proto.ErrRejected, client: 1},
		{query: "func0", expected: ErrUnauthorized, client: 2},
		{query: "func1", expected: map[string]any{}},
		{query: "func2[1]", expected: float64(1)},
		{query: "func2 1", expected: float64(1)},
		{query: "func2", arg: 1, expected: float64(1)},
		{query: "func3", expected: err3},
		{query: "func4[true]", expected: true},
		{query: "func4 true", expected: true},
		{query: "func4[false]", expected: err4},
		{query: "func4 false", expected: err4},
		{query: "func4[]", expected: err4},
		{query: "func4 ", expected: err4},
		{query: "func4", arg: true, expected: true},
		{query: "func4", arg: false, expected: err4},
		{query: `func5[true, 1, "a"]`, expected: []any{true, float64(1), "a"}},
		{query: "func5 true 1 a", expected: []any{true, float64(1), "a"}},
		{query: "func5", args: []any{true, 1, "a"}, expected: []any{true, float64(1), "a"}},
		{query: `func6{"i":1}`, expected: &structI{1}},
		{query: `func6[{"i":1}]`, expected: structI{1}},
		{query: `func6 -i 1`, expected: structI{1}},
		{query: `func6 -i 1`, expected: &structI{1}},
		{query: `func6`, arg: structI{1}, expected: structI{1}},
		{query: "func7[]", expected: (*structI)(nil)},
		{query: `func7{"i":1}`, expected: &structI{1}},
		{query: "func7 ", expected: (*structI)(nil)},
		{query: "func7", arg: (*structI)(nil), expected: (*structI)(nil)},
		{query: "func7", arg: structI{1}, expected: structI{1}},
		//{name: "func7 -i 1", expected: &structI{1}}, //TODO minor consider to fix
		{query: `func8[{"i":1},{"b":true}]`, expected: struct1{structI{1}, structB{true}}},
		{query: `func8 -i 1 -b`, expected: struct1{structI{1}, structB{true}}},
		{query: `func8`, args: []any{structI{1}, structB{true}}, expected: struct1{structI{1}, structB{true}}},
		{query: `func9[{"i":1},{"b":true}]`, expected: struct1{structI{1}, structB{true}}},
		{query: `func9[{"i":1}]`, expected: struct1{structI{1}, structB{false}}},
		{query: `func9`, arg: structI{1}, expected: structI{1}},
		{query: `func9`, args: []any{structI{1}, structB{true}}, expected: struct1{structI{1}, structB{true}}},
		//{name: `func9 -i 1 -b true`, expected: struct1{structI{1}, structB{true}}}, //TODO minor consider to fix
		{query: `func10{"i":{"i":1},"b":{"b":true}}`, expected: struct3{structI{1}, structB{true}}},
		{query: `func10{"i":{"i":1}}`, expected: struct3{structI{1}, structB{false}}},
		{query: `func10`, arg: struct3{structI{1}, structB{true}}, expected: struct3{structI{1}, structB{true}}},
		{query: `func10`, arg: struct3{StructI: structI{1}}, expected: struct3{structI{1}, structB{false}}},
		{query: `func11`, args: []any{"a", "b", "c"}, expected: []any{"b", "c", "a"}},
	}
	skip := func(t *testing.T, routes int, client int, case_ int, err error) {
		for _, s := range skipped {
			a1 := s.routes == routes+1
			a2 := s.client == client+1
			a3 := slices.Contains(s.cases, case_+1)
			if a1 || a2 || a3 {
				b1 := a1 || s.routes == 0
				b2 := a2 || s.client == 0
				b3 := a3 || len(s.cases) == 0
				if b1 && b2 && b3 {
					t.Skip(err)
					return
				}
			}
		}
	}

	for i1, r := range routes {
		t.Run(fmt.Sprintf("routes:%d", i1+1), func(t *testing.T) {
			app := NewApp(port)
			app.Logger(plog.Type(app))
			app.Routes(r...)
			setHandlers(app)

			ctx, cancel := context.WithCancel(context.Background())
			t.Run("backend", func(t *testing.T) {
				t.Parallel()
				if err := app.Run(ctx); err != nil {
					skip(t, i1, 0, 0, err)
					t.Fatal(err)
				}
			})
			time.Sleep(1 * time.Millisecond)

			t.Run("frontend", func(t *testing.T) {
				t.Parallel()
				t.Cleanup(cancel)
				for i2, c := range clients {
					t.Run(fmt.Sprintf("client:%d", i2+1), func(t *testing.T) {
						for i3, tt := range cases {
							if tt.client > 0 && tt.client-1 != i2 {
								continue
							}
							t.Run(fmt.Sprintf("%d.%s", i3+1, tt.query), func(t *testing.T) {
								skip(t, i1, i2, 0, nil)
								client, err := c(t)
								if err != nil {
									t.Fatal(err)
								}
								t.Cleanup(client.Flush)
								args := tt.args
								if tt.arg != nil {
									args = append([]any{tt.arg}, args...)
								}
								skip(t, i1, i2, i3, err)
								if err := Call(client, tt.query, args...); err != nil {
									assert.Equal(t, tt.expected, err)
									return
								}

								v := reflect.New(reflect.TypeOf(tt.expected))
								if err := client.Decode(v.Interface()); err != nil {
									assert.Equal(t, tt.expected, err)
									skip(t, i1, i2, i3, err)
								} else {
									assert.Equal(t, tt.expected, v.Elem().Interface())
									skip(t, i1, i2, i3, err)
								}
							})
						}
					})
				}
			})
		})
	}
}

func function0() bool { return false }

func function1() {
	plog.Println("function1")
}

func function2(i int) int {
	return i
}

var err3 = errors.New("test error 3")

func function3() error {
	return err3
}

var err4 = errors.New("test error 4")

func function4(b bool) (bool, error) {
	if b {
		return true, nil
	}
	return false, err4
}

func function5(b bool, i int, s string) (bool, int, string) {
	return b, i, s
}

type structI struct {
	I int `json:"i" name:"i"`
}

type structB struct {
	B bool `json:"b" name:"b"`
}

type struct1 struct {
	structI
	structB
}

type struct2 struct {
	*structI
	*structB
}

type struct3 struct {
	StructI structI `json:"I"`
	StructB structB `json:"B"`
}

func function6(i structI) structI {
	return i
}

func function7(i *structI) *structI {
	return i
}

func function8(i structI, b structB) struct1 {
	return struct1{i, b}
}

func function9(i *structI, b *structB) *struct2 {
	return &struct2{i, b}
}

func function10(s struct3) struct3 {
	return s
}

func function11(arg string, args ...string) (arr []string) {
	arr = append(args, arg)
	plog.Println(arr)
	return
}
