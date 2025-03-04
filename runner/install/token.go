package install

import (
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/client/apphost"
	"github.com/cryptopunkscc/portal/client/keys"
	"github.com/cryptopunkscc/portal/runtime/apps"
)

func (i Install) Token(pkg string) (token mod.AccessToken, err error) {
	tokens := apps.Tokens{Dir: i.appsDir}
	client := apphost.NewClient()

	if token, err = tokens.GetToken(pkg); err == nil {
		return
	}

	t := &mod.AccessToken{}
	if t.Identity, err = api.DefaultClient.Resolve(pkg); err == nil {
		if at, err2 := client.ListTokens(nil); err2 == nil {
			for _, tt := range at {
				if tt.Identity.IsEqual(t.Identity) {
					token = tt
					return
				}
			}
		}

	} else if t.Identity, err = keys.NewClient().CreateKey(pkg); err != nil {
		return
	}

	args := apphost.CreateTokenArgs{ID: t.Identity}
	if t, err = client.CreateToken(args); err != nil {
		return
	}

	token = *t
	err = tokens.SetToken(pkg, token)
	return
}
