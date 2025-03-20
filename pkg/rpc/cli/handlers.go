package cli

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/pkg/rpc/caller/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/caller/query"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/rpc/router"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"log"
)

var Handler = cmd.Handler{
	Name: "cli",
	Params: cmd.Params{
		{
			Name: "interactive i",
			Type: "bool",
			Desc: "Enable interactive mode.",
		},
	},
	Func: func(r *router.Base, c *stream.Client, opt optCli) {
		c.Ending = nil
		c.Marshal = Marshal
		r.Unmarshal = cli.Unmarshal
		if !opt.Interactive {
			r.Responses = 1
		}
	},
}

type optCli struct {
	Interactive bool `cli:"interactive i" query:"interactive i"`
}

type optCoding struct {
	Format  string `cli:"format f" query:"format"`
	Encoder string `cli:"encoder e" query:"encoder"`
	Decoder string `cli:"decoder d" query:"decoder"`
	Limit   int64  `cli:"limit l" query:"limit"`
}

var EncodingHandler = cmd.Handler{
	Name: "encoding",
	Params: cmd.Params{
		{
			Name: "encoder e",
			Type: "string",
			Desc: "Encoder [json, cli]",
		},
		{
			Name: "decoder d",
			Type: "string",
			Desc: "Decoder [json]",
		},
		{
			Name: "format f",
			Type: "string",
			Desc: "Format [query, cli]",
		},
	},
	Func: func(r *router.Base, c *stream.Client, opt optCoding) {
		log.Printf("encoding encoder: %s, decoder: %s, format: %s", opt.Encoder, opt.Decoder, opt.Format)
		log.Printf("%v %v %v", c.Marshal, c.Unmarshal, r.Unmarshal)
		switch opt.Encoder {
		case "json":
			c.Marshal = json.Marshal
			c.Ending = []byte("\n")
		case "cli":
			c.Marshal = Marshal
		}
		switch opt.Decoder {
		case "json":
			c.Unmarshal = json.Unmarshal
		}
		switch opt.Format {
		case "query":
			r.Unmarshal = query.Unmarshal
		case "cli":
			r.Unmarshal = cli.Unmarshal
		}
		if opt.Limit != 0 {
			r.Responses = opt.Limit
		}
		log.Printf("%v %v %v", c.Marshal, c.Unmarshal, r.Unmarshal)
	},
}

var InteractiveModeHandlers = cmd.Handlers{
	{
		Name: "-i",
		Desc: "Run interactive mode",
		Func: func(runner *Runner) { runner.interactive = true },
	},
	{
		Name: "exit",
		Desc: "Exit interactive mode",
		Func: func(runner *Runner) { runner.interactive = false },
	},
}
