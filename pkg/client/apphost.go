package client

import (
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/apphost"
	apphost2 "github.com/cryptopunkscc/astrald/mod/apphost/client"
)

type Apphost struct {
	astrald *astrald.Client
	apphost2.Client
}

func (a *Apphost) ListTokens(ctx *astral.Context, out string) ([]apphost.AccessToken, error) {
	args := query.Args{}
	if out != "" {
		args = query.Args{"out": out}
	}
	return List[apphost.AccessToken](ctx, *a.astrald, "apphost.list_tokens", args)
}

type AccessTokens []apphost.AccessToken

func (t AccessTokens) MarshalCLI() (s string) {
	for i, token := range t {
		s += fmt.Sprintf("%d: %s %s\n", i, token.Identity, token.Token)
	}
	return
}

func (a *Apphost) SignAppContract(ctx *astral.Context, id *astral.Identity) (out *astral.ObjectID, err error) {
	return Receive[*astral.ObjectID](ctx, *a.astrald, "apphost.sign_app_contract", query.Args{"id": id.String()})
}
