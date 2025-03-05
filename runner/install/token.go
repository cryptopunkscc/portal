package install

import (
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/client/apphost"
	"github.com/cryptopunkscc/portal/client/keys"
	"github.com/cryptopunkscc/portal/runtime/tokens"
)

func (i Install) Token(pkg string) (token mod.AccessToken, err error) {
	repo := tokens.Repository{}
	client := apphost.NewClient()

	if token, err = repo.Get(pkg); err == nil {
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
	err = repo.Set(pkg, token)
	return
}
