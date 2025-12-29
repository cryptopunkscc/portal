package apphost

import (
	"errors"
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Op(client Client, target ...string) (c OpClient) { return OpClient{client, Target(target...)} }

type OpClient struct {
	Client
	Target string
}

func (c OpClient) r() rpc.Conn {
	return c.Rpc().Request(c.Target, "apphost")
}

type CreateTokenArgs struct {
	ID       *astral.Identity `query:"id" cli:"id"`
	Duration *astral.Duration `query:"duration" cli:"duration d"`
	Out      astral.String    `query:"out" cli:"out o"`
}

func (c OpClient) CreateToken(args CreateTokenArgs) (ac *mod.AccessToken, err error) {
	if args.ID == nil {
		return nil, errors.New("id is required")
	}
	if args.Out == "" {
		args.Out = "json"
	}
	r, err := rpc.Query[rpc.Json[*mod.AccessToken]](c.r(), "create_token", args)
	if err != nil {
		return
	}
	ac = r.Object
	return
}

type ListTokensArgs struct {
	ID  *astral.Identity `query:"id" cli:"id"`
	Out string           `query:"format" cli:"format f"`
}

func (c OpClient) ListTokens(args *ListTokensArgs) (at AccessTokens, err error) {
	if args == nil {
		args = &ListTokensArgs{}
	}
	if args.Out == "" {
		args.Out = "json"
	}
	result, err := rpc.Subscribe[rpc.Json[*mod.AccessToken]](c.r(), "list_tokens", rpc.Opt{"out": "json"})
	if err != nil {
		return
	}
	for r := range result {
		at = append(at, *r.Object)
	}
	return
}

type AccessTokens []mod.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Identity, token.Token)
	}
	return
}

func (c OpClient) SignAppContract(id *astral.Identity) (out *astral.ObjectID, err error) {
	s, err := rpc.Query[rpc.Json[*astral.ObjectID]](c.r(), "sign_app_contract", rpc.Opt{
		"id":  id.String(),
		"out": "json",
	})
	if err != nil {
		return
	}
	out = s.Object
	return
}
