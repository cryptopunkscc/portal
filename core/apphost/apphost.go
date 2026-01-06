package apphost

import (
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/apphost"
)

func (a *Adapter) CreateToken(id *astral.Identity) (ac *apphost.AccessToken, err error) {
	return Request[*apphost.AccessToken](nil, *a.Client, "apphost.create_token", query.Args{"id": id.String()})
}

func (a *Adapter) ListTokens(out string) ([]apphost.AccessToken, error) {
	args := query.Args{}
	if out != "" {
		args = query.Args{"out": out}
	}
	return List[apphost.AccessToken](nil, *a.Client, "apphost.list_tokens", args)
}

type AccessTokens []apphost.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Identity, token.Token)
	}
	return
}

// SignAppContract signs contract with given app and returns contract identity
func (a *Adapter) SignAppContract(id *astral.Identity) (out *astral.ObjectID, err error) {
	return Request[*astral.ObjectID](nil, *a.Client, "apphost.sign_app_contract", query.Args{"id": id.String()})
}
