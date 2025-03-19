package token

import (
	"github.com/cryptopunkscc/astrald/astral"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/dir"
)

type Repository struct {
	Dir string
	*apphost.Adapter
}

func (r Repository) Set(pkg string, token *mod.AccessToken) (err error) {
	if r.Dir == "" {
		r.Dir = dir.Token
	}
	if err = os.WriteJson[*mod.AccessToken](token, r.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (r Repository) Get(pkg string) (token *mod.AccessToken, err error) {
	if r.Dir == "" {
		r.Dir = dir.Token
	}
	if token, err = os.ReadJson[*mod.AccessToken](r.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (r Repository) Fetch(pkg string) (accessToken *mod.AccessToken, err error) {
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	if accessToken, err = r.Get(pkg); err == nil {
		return
	}

	i, err := r.Adapter.Resolve(pkg)
	if err != nil {
		return
	}

	at, err := r.Adapter.Token().List(nil)
	if err != nil {
		return
	}

	for _, t := range at {
		if t.Identity.IsEqual(i) {
			accessToken = &t
			err = r.Set(pkg, accessToken)
			return
		}
	}
	return
}

func (r Repository) Resolve(pkg string) (accessToken *mod.AccessToken, err error) {
	if r.Adapter == nil {
		r.Adapter = apphost.Default
	}
	if accessToken, err = r.Fetch(pkg); err == nil {
		return
	}

	var i *astral.Identity
	if i, err = r.Adapter.Key().Create(pkg); err != nil {
		return
	}

	args := apphost.CreateTokenArgs{ID: i}
	if accessToken, err = r.Adapter.Token().Create(args); err != nil {
		return
	}

	err = r.Set(pkg, accessToken)
	return
}
