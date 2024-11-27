package main

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
	"math"
	"time"
)

func main() {
	ctx := context.Background()

	handler := cmd.Handler{
		Name: "test-cli",
		Desc: "Test cli",
		Func: apphost.ServeFunc,
		Sub: cmd.Handlers{
			{Func: add, Name: "sum s", Desc: "Sum given values.", Params: cmd.Params{
				{Type: "int", Desc: "a."},
				{Type: "int", Desc: "b."},
			}},
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
			{Func: ticker, Name: "ticker t", Desc: "Start ticker.", Params: cmd.Params{
				{Type: "int", Desc: "Counter limit."},
			}},
			{
				Name: "echo e",
				Desc: "Echo command.",
				Params: cmd.Params{
					{Type: "string"},
				},
				Func: func(str string) string { return str },
			},
			{
				Name: "test",
				Desc: "Test command.",
				Params: append(
					Options{}.CmdParams(),
					cmd.Param{Type: "string"},
				),
				Func: func(o Options, str string) string {
					log.Println(o)
					return fmt.Sprintf("%v %d %s %s", o.B, o.I, o.S, str)
				},
			},
			apphost.Serve(),
		},
	}

	if err := cli.New(handler).Run(ctx); err != nil {
		panic(err)
	}
}

func ticker(amount int) <-chan int {
	ch := make(chan int)
	if amount == 0 {
		amount = math.MaxInt
	}
	go func() {
		defer close(ch)
		for i := 1; i < amount+1; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
	}()
	return ch
}

func add(a int, b int) int {
	return a + b
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

type Options struct {
	S string `cli:"s" json:"s"`
	I int    `cli:"i" json:"i"`
	B bool   `cli:"b" json:"b"`
}

func (o Options) CmdParams() cmd.Params {
	return cmd.Params{
		{Name: "s", Type: "string"},
		{Name: "i", Type: "int"},
		{Name: "b", Type: "bool"},
	}
}
