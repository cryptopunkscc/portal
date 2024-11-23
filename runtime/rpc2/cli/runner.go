package cli

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/json"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/cryptopunkscc/portal/runtime/rpc2/router"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"io"
)

type Runner struct {
	router.Base
	conn        stream.Serializer
	interactive bool
}

func New(handler cmd.Handler) (runner *Runner) {
	root := cmd.Root(handler)

	injectHelp(&handler)

	runner = &Runner{
		conn: cliConnection(),
	}

	r := router.Base{
		Unmarshalers: []caller.Unmarshaler{
			cli.Unmarshaler{},
			json.Unmarshaler{},
			query.Unmarshaler{},
		},
	}
	r.Dependencies = []any{&root, &r, runner}
	r.Registry = router.CreateRegistry(handler)

	runner.Base = r

	return
}

func InteractiveMode() cmd.Handlers {
	return cmd.Handlers{
		{Name: "-i", Desc: "Run interactive mode", Func: func(runner *Runner) { runner.interactive = true }},
		{Name: "exit", Desc: "Exit interactive mode", Func: func(runner *Runner) { runner.interactive = false }},
	}
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
