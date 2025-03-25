package astrald

import (
	"github.com/cryptopunkscc/portal/core/token"
)

func (r *Initializer) fetchAuthToken(pkg string) (err error) {
	t, err := r.tokens().Get(pkg)
	if err != nil {
		return
	}
	r.Apphost.AuthToken = string(t.Token)
	r.log.Println("fetched", pkg, "auth token")
	return
}

func (r *Initializer) resolveAuthToken(pkg string) (err error) {
	t, err := r.tokens().Resolve(pkg)
	if err != nil {
		return
	}
	r.Apphost.AuthToken = string(t.Token)
	r.log.Println("resolved", pkg, "auth token")
	return
}

func (r *Initializer) tokens() *token.Repository {
	return token.NewRepository(r.TokensDir, r.Apphost)
}
