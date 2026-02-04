package apphost

import (
	"errors"
	"os"

	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/env"
	pkgOs "github.com/cryptopunkscc/portal/pkg/util/os"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

func (a *Adapter) Tokens(dir string) *Tokens {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	return &Tokens{Dir: dir, Adapter: a}
}

type Tokens struct {
	Dir string
	*Adapter
}

func (r *Tokens) Load() bool {
	if r.Dir == "" && env.PortaldTokens.Exist() {
		r.Dir = env.PortaldTokens.MkdirAll()
	}
	return len(r.Dir) > 0
}

func (r *Tokens) dir() string {
	r.Load()
	if len(r.Dir) == 0 {
		panic("no tokens dir")
	}
	return r.Dir
}

func (r *Tokens) Set(pkg string, token *mod.AccessToken) (err error) {
	if err = pkgOs.WriteJson[*mod.AccessToken](token, r.dir(), pkg); err != nil {
		err = plog.Err(err, pkg)
	}
	return
}

func (r *Tokens) Get(pkg string) (token *mod.AccessToken, err error) {
	if token, err = pkgOs.ReadJson[*mod.AccessToken](r.dir(), pkg); err != nil {
		err = plog.Err(ErrNotCached, pkg)
	}
	return
}

var ErrNotCached = errors.New("apphost auth token is not cached or cannot be loaded")

func (r *Tokens) Resolve(pkg string) (accessToken *mod.AccessToken, err error) {
	defer plog.TraceErr(&err)
	if accessToken, err = r.Get(pkg); err == nil {
		return
	}

	id, _ := r.Adapter.Resolve(pkg)

	if id != nil {
		var tokens []mod.AccessToken
		if tokens, err = r.Adapter.ListTokens(""); err != nil {
			return
		}

		for _, t := range tokens {
			if t.Identity.IsEqual(id) {
				accessToken = &t
				err = r.Set(pkg, accessToken)
				return
			}
		}
	} else if id, err = r.Adapter.Keys().CreateKey(nil, pkg); err != nil {
		return
	}

	if accessToken, err = r.Adapter.CreateToken(id); err != nil {
		return
	}

	if err = r.Set(pkg, accessToken); err != nil {
		return
	}
	return
}
