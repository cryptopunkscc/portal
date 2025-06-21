package token

import (
	"errors"
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

func (r *Repository) dir() string {
	if r.Dir == "" {
		r.Dir = env.PortaldTokens.MkdirAll()
	}
	return r.Dir
}

func (r *Repository) apphost() api.Client {
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	return r.Adapter
}

func (r *Repository) op() api.OpClient {
	return api.Op(r.apphost())
}

func (r *Repository) Set(pkg string, token *mod.AccessToken) (err error) {
	if err = pkgOs.WriteJson[*mod.AccessToken](token, r.dir(), pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (r *Repository) Get(pkg string) (token *mod.AccessToken, err error) {
	if token, err = pkgOs.ReadJson[*mod.AccessToken](r.dir(), pkg); err != nil {
		err = plog.Err(ErrNotCached)
	}
	return
}

var ErrNotCached = errors.New("apphost auth token is not cached or cannot be loaded")

func (r *Repository) List(args *api.ListTokensArgs) (api.AccessTokens, error) {
	return r.op().ListTokens(args)
}

func (r *Repository) Resolve(pkg string) (accessToken *mod.AccessToken, err error) {
	defer plog.TraceErr(&err)
	if accessToken, err = r.Get(pkg); err == nil {
		return
	}

	id, _ := r.apphost().Resolve(pkg)

	if id != nil {
		var tokens []mod.AccessToken
		if tokens, err = r.op().ListTokens(nil); err != nil {
			return
		}

		for _, t := range tokens {
			if t.Identity.IsEqual(id) {
				accessToken = &t
				err = r.Set(pkg, accessToken)
				return
			}
		}
	} else if id, err = keys.Op(r.apphost()).CreateKey(pkg); err != nil {
		return
	}

	args := api.CreateTokenArgs{ID: id}
	if accessToken, err = r.op().CreateToken(args); err != nil {
		return
	}

	if err = r.Set(pkg, accessToken); err != nil {
		return
	}
	return
}
