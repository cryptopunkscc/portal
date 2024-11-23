package main

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	ctx := context.Background()

	handler := cmd.Handler{
		Name: "test-cli",
		Desc: "Test cli",
		Sub: cmd.Handlers{
			{Func: inc, Name: "inc i", Desc: "Increment given value.", Params: cmd.Params{
				{Type: "int", Desc: "value to decrement."},
			}},
			{Func: dec, Name: "dec d", Desc: "Decrement given value.", Params: cmd.Params{
				{Type: "int", Desc: "value to increment."},
			}},
			{Func: foo, Name: "foo f", Sub: cmd.Handlers{
				{Func: bar, Name: "bar b"},
				{Func: baz, Name: "baz"},
			}},
		},
	}

	if err := cli.New(handler).Run(ctx); err != nil {
		panic(err)
	}
}

func inc(a int) int {
	return a + 1
}

func dec(a int) int {
	return a - 1
}

func foo() string {
	return "foo"
}

func bar() string {
	return "bar"
}

func baz() string {
	return "baz"
}
