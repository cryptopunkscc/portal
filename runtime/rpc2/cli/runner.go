package cli

import (
	"context"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/clir"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/json"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"io"
)

type Runner struct {
	rpc.Router
	conn        rpc.Serializer
	interactive bool
}

func New(handler cmd.Handler) (runner *Runner) {

	handler.AddSub(cmd.Handlers{
		{Name: "-i", Desc: "Run interactive mode", Func: func() { runner.interactive = true }},
		{Name: "exit", Desc: "Exit interactive mode", Func: func() { runner.interactive = false }},
	}...)

	injectHelp(&handler)
	router := rpc.Router{
		Registry: rpc.CreateRegistry(handler),
		Unmarshalers: []caller.Unmarshaler{
			json.Unmarshaler{},
			clir.Unmarshaler{},
		},
	}

	runner = &Runner{
		Router: router,
		conn:   cliConnection(),
	}

	return
}

func (c *Runner) Run(_ context.Context) error {
	for {
		// read query
		bytes, err := c.conn.Bytes()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		query := string(bytes)

		// handle query
		result := c.Query(query).Call()
		if result != nil {
			err = c.conn.Encode(result)
			if err != nil {
				return err
			}
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
