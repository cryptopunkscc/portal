package token

import (
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/keys"
	"github.com/cryptopunkscc/portal/core/apphost"
	pkgOs "github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
)

type Repository struct {
	Dir string
	*apphost.Adapter
}

func NewRepository(
	dir string,
	adapter *apphost.Adapter,
) *Repository {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	return &Repository{Dir: dir, Adapter: adapter}
}

func (r *Repository) Set(pkg string, token *mod.AccessToken) (err error) {
	if r.Dir == "" {
		r.Dir = env.PortaldTokens.MkdirAll()
	}
	if err = pkgOs.WriteJson[*mod.AccessToken](token, r.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (r *Repository) Get(pkg string) (token *mod.AccessToken, err error) {
	if r.Dir == "" {
		r.Dir = env.PortaldTokens.MkdirAll()
	}
	if token, err = pkgOs.ReadJson[*mod.AccessToken](r.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (r *Repository) List(args *api.ListTokensArgs) (api.AccessTokens, error) {
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	return api.TokenClient(r).List(args)
}

func (r *Repository) Resolve(pkg string) (accessToken *mod.AccessToken, err error) {
	defer plog.TraceErr(&err)
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	if accessToken, err = r.Get(pkg); err == nil {
		return
	}

	id, _ := r.Adapter.Resolve(pkg)

	if id != nil {
		var tokens []mod.AccessToken
		if tokens, err = api.TokenClient(r).List(nil); err != nil {
			return
		}

		for _, t := range tokens {
			if t.Identity.IsEqual(id) {
				accessToken = &t
				err = r.Set(pkg, accessToken)
				return
			}
		}
	} else if id, err = keys.Client(r).Create(pkg); err != nil {
		return
	}

	args := api.CreateTokenArgs{ID: id}
	if accessToken, err = api.TokenClient(r).Create(args); err != nil {
		return
	}

	if err = r.Set(pkg, accessToken); err != nil {
		return
	}
	return
}
