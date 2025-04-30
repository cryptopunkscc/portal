package apphost

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func TokenClient(rpc rpc.Rpc) TokenConn {
	return TokenConn{rpc.Request("localnode")}
}

type TokenConn struct{ rpc.Conn }

type CreateTokenArgs struct {
	ID       *astral.Identity `query:"id" cli:"id"`
	Duration *astral.Duration `query:"duration" cli:"duration d"`
	Out      astral.String    `query:"out" cli:"out o"`
}

func (c TokenConn) Create(args CreateTokenArgs) (ac *mod.AccessToken, err error) {
	if args.ID == nil {
		return nil, errors.New("id is required")
	}
	if args.Out == "" {
		args.Out = "json"
	}
	r, err := rpc.Query[rpc.JsonObject[*mod.AccessToken]](c, "apphost.create_token", args)
	if err != nil {
		return
	}
	ac = r.Payload
	return
}

type ListTokensArgs struct {
	ID     *astral.Identity `query:"id" cli:"id"`
	Format string           `query:"format" cli:"format f"`
}

func (c TokenConn) List(args *ListTokensArgs) (AccessTokens, error) {
	if args == nil {
		args = &ListTokensArgs{}
	}
	if args.Format == "" {
		args.Format = "json"
	}
	return rpc.Query[AccessTokens](c, "apphost.list_tokens", args)
}

type AccessTokens []mod.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Token, token.Identity)
	}
	return
}
