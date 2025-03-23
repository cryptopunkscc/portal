package astrald

import (
	"github.com/cryptopunkscc/portal/core/token"
)

func (r *Runner) fetchAuthToken(pkg string) (err error) {
	t, err := r.tokens().Get(pkg)
	if err != nil {
		return
	}
	r.ResolvedTokens.Set(pkg, t)
	r.log.Println("fetched", pkg, "auth token")
	return
}

func (r *Runner) resolveAuthToken(pkg string) (err error) {
	t, err := r.tokens().Resolve(pkg)
	if err != nil {
		return
	}
	r.ResolvedTokens.Set(pkg, t)
	r.log.Println("resolved", pkg, "auth token")
	return
}

func (r *Runner) tokens() *token.Repository {
	return token.NewRepository(r.TokensDir, r.Apphost)
}
