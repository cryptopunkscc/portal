package astrald

import (
	"github.com/cryptopunkscc/portal/core/token"
)

func (i *Initializer) fetchAuthToken(pkg string) (err error) {
	t, err := i.tokens().Get(pkg)
	if err != nil {
		return
	}
	i.Apphost.AuthToken = string(t.Token)
	i.log.Println("fetched", pkg, "auth token")
	return
}

func (i *Initializer) resolveAuthToken(pkg string) (err error) {
	t, err := i.tokens().Resolve(pkg)
	if err != nil {
		return
	}
	i.Apphost.AuthToken = string(t.Token)
	i.log.Println("resolved", pkg, "auth token")
	return
}

func (i *Initializer) tokens() *token.Repository {
	return token.NewRepository(i.TokensDir, i.Apphost)
}
