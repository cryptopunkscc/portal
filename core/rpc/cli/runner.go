package cli

import (
	"context"
	"github.com/cryptopunkscc/portal/core/rpc/caller/cli"
	"github.com/cryptopunkscc/portal/core/rpc/cmd"
	"github.com/cryptopunkscc/portal/core/rpc/router"
	"github.com/cryptopunkscc/portal/core/rpc/stream"
	"io"
)

type Runner struct {
	router.Base
	Conn        stream.Serializer
	interactive bool
}

func New(handler cmd.Handler) (runner *Runner) {
	root := cmd.Root(handler)

	runner = &Runner{
		Conn: cliConnection(),
	}

	r := router.Base{
		Unmarshal: cli.Unmarshal,
	}
	r.Dependencies = []any{&root, &r, runner}
	r.Registry = router.CreateRegistry(handler)

	runner.Base = r

	return
}

func (c *Runner) Run(ctx context.Context) error {
	for {
		// read query
		bytes, err := c.Conn.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// handle query
		rr := *c
		rr.Setup(bytes)
		rr.Dependencies = append([]any{ctx}, rr.Dependencies...)
		if err = rr.Respond(&c.Conn); err != nil {
			return err
		}

		// interactive mode check
		if !c.interactive {
			return nil
		}
		if _, err = c.Conn.Write([]byte("$ ")); err != nil {
			return err
		}
	}
}
