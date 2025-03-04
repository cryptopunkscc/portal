package apphost

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	apphost2 "github.com/cryptopunkscc/astrald/mod/apphost"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
)

func NewClient() Client { return Client{apphost.Request("localnode")} }

type Client struct{ rpc.Conn }

type CreateTokenArgs struct {
	ID       *astral.Identity `query:"id" cli:"id"`
	Duration *astral.Duration `query:"duration" cli:"duration d"`
	Format   astral.String    `query:"format" cli:"format f"`
}

func (c Client) CreateToken(args CreateTokenArgs) (*apphost2.AccessToken, error) {
	if args.ID == nil {
		return nil, errors.New("id is required")
	}
	if args.Format == "" {
		args.Format = "json"
	}
	return rpc.Query[*apphost2.AccessToken](c, "apphost.create_token", args)
}

type ListTokensArgs struct {
	ID     *astral.Identity `query:"id" cli:"id"`
	Format string           `query:"format" cli:"format f"`
}

func (c Client) ListTokens(args ListTokensArgs) (AccessTokens, error) {
	if args.Format == "" {
		args.Format = "json"
	}
	return rpc.Query[AccessTokens](c, "apphost.list_tokens", args)
}

type AccessTokens []apphost2.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Token, token.Identity)
	}
	return
}
