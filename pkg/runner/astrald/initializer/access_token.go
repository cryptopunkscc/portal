package initializer

import (
	"github.com/cryptopunkscc/portal/pkg/client"
)

func (i *Astrald) fetchAuthToken(pkg string) (err error) {
	t, err := i.tokens().Get(pkg)
	if err != nil {
		return
	}
	i.Client.Token = string(t.Token)
	i.log.Println("fetched", pkg, "auth token")
	return
}

func (i *Astrald) resolveAuthToken(pkg string) (err error) {
	t, err := i.tokens().Resolve(pkg)
	if err != nil {
		return
	}
	i.Client.Token = string(t.Token)
	i.log.Println("resolved", pkg, "auth token")
	return
}

func (i *Astrald) tokens() *client.Tokens {
	return i.Client.Tokens(i.TokensDir)
}
