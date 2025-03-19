package setup

import (
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

func (r *Runner) resolveAuthToken(pkg string) (err error) {
	if r.Tokens.Adapter == nil {
		r.Tokens.Adapter = &apphost.Adapter{}
		r.Tokens.AuthToken = r.nodeAuthToken
	}
	for _, endpoint := range r.ApphostConfig.Listen[:1] {
		r.Tokens.Endpoint = endpoint
	}
	r.log.D().Println("using auth token: [", r.nodeAuthToken, "] and apphost address:", r.Tokens.Endpoint)
	t, err := r.Tokens.Resolve(pkg)
	if err != nil {
		return
	}
	r.ResolvedTokens.Set(pkg, t)
	r.log.Println("resolved", pkg, "auth token")
	return
}
