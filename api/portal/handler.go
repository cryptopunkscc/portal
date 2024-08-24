package portal

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

type Service interface {
	Query(ctx context.Context, query string) error
}

func Handler(service Service) cmd.Handler {
	return cmd.Handler{
		Func: service.Query,
		Name: "portal",
		Desc: "Portal command line.",
		Params: cmd.Params{{
			Type: "string",
			Desc: "Portal app query. Accepted formats are CLI or URL with query or JSON args",
		}},
	}
}
