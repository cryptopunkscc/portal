package apphost

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func TokenClient(rpc rpc.Rpc) TokenConn {
	return TokenConn{rpc.Request("localnode", "apphost")}
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
	r, err := rpc.Query[rpc.Json[*mod.AccessToken]](c, "create_token", args)
	if err != nil {
		return
	}
	ac = r.Object
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
	return rpc.Query[AccessTokens](c, "list_tokens", args)
}

type AccessTokens []mod.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Token, token.Identity)
	}
	return
}

func (c TokenConn) SignAppContract(id *astral.Identity) (out *astral.ObjectID, err error) {
	args := opSignAppContractArgs{
		ID:  id,
		Out: "json",
	}
	s, err := rpc.Query[rpc.Json[*astral.ObjectID]](c, "sign_app_contract", args)
	if err != nil {
		return
	}
	out = s.Object
	return
}

type opSignAppContractArgs struct {
	ID       *astral.Identity
	Out      string           `query:"out"`
	Duration *astral.Duration `query:"duration"`
}
