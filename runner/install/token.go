package install

import (
	apphost2 "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/client/apphost"
	"github.com/cryptopunkscc/portal/client/keys"
)

func Token(alias string) (token *apphost2.AccessToken, err error) {
	key, err := keys.NewClient().CreateKey(alias)
	if err != nil {
		return
	}
	return apphost.NewClient().CreateToken(apphost.CreateTokenArgs{ID: key})
}
