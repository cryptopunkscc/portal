package token

import (
	"errors"
	"os"

	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/keys"
	"github.com/cryptopunkscc/portal/core/apphost"
	pkgOs "github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/plog"
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

func (r *Repository) Load() bool {
	if r.Dir == "" && env.PortaldTokens.Exist() {
		r.Dir = env.PortaldTokens.MkdirAll()
	}
	return len(r.Dir) > 0
}

func (r *Repository) dir() string {
	r.Load()
	if len(r.Dir) == 0 {
		panic("no tokens dir")
	}
	return r.Dir
}

func (r *Repository) apphost() api.Client {
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	return r.Adapter
}

func (r *Repository) Set(pkg string, token *mod.AccessToken) (err error) {
	if err = pkgOs.WriteJson[*mod.AccessToken](token, r.dir(), pkg); err != nil {
		err = plog.Err(err, pkg)
	}
	return
}

func (r *Repository) Get(pkg string) (token *mod.AccessToken, err error) {
	if token, err = pkgOs.ReadJson[*mod.AccessToken](r.dir(), pkg); err != nil {
		err = plog.Err(ErrNotCached, pkg)
	}
	return
}

var ErrNotCached = errors.New("apphost auth token is not cached or cannot be loaded")

type ListTokensArgs struct {
	Out string `query:"format" cli:"format f"`
}
type AccessTokens []mod.AccessToken

func (r *Repository) List(args ListTokensArgs) (AccessTokens, error) {
	return r.Adapter.ListTokens(args.Out)
}

func (r *Repository) Resolve(pkg string) (accessToken *mod.AccessToken, err error) {
	defer plog.TraceErr(&err)
	if accessToken, err = r.Get(pkg); err == nil {
		return
	}

	id, _ := r.apphost().Resolve(pkg)

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
	} else if id, err = keys.Op(r.apphost()).CreateKey(pkg); err != nil {
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
