package cli

import (
	"context"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/json"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"io"
)

type Runner struct {
	rpc.Router
	conn        rpc.Serializer
	interactive bool
}

func New(handler cmd.Handler) (runner *Runner) {
	root := cmd.Root(handler)

	handler.AddSub(cmd.Handlers{
		{Name: "-i", Desc: "Run interactive mode", Func: func() { runner.interactive = true }},
		{Name: "exit", Desc: "Exit interactive mode", Func: func() { runner.interactive = false }},
	}...)

	injectHelp(&handler)
	router := rpc.Router{
		Unmarshalers: []caller.Unmarshaler{
			cli.Unmarshaler{},
			json.Unmarshaler{},
			query.Unmarshaler{},
		},
	}
	router.Dependencies = []any{&root, &router}
	router.Registry = rpc.CreateRegistry(handler)

	runner = &Runner{
		Router: router,
		conn:   cliConnection(),
	}

	return
}

func (c *Runner) Run(ctx context.Context) error {
	for {
		// read query
		bytes, err := c.conn.Bytes()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// handle query
		rr := c.Query(string(bytes))
		rr.Dependencies = append([]any{ctx}, rr.Dependencies...)
		if err = rr.Respond(&c.conn); err != nil {
			return err
		}

		// interactive mode check
		if !c.interactive {
			return nil
		}
		if _, err = c.conn.Write([]byte("$ ")); err != nil {
			return err
		}
	}
}
