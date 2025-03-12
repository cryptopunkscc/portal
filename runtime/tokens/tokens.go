package tokens

import (
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/dir"
)

type Repository struct {
	Dir string
}

func (t Repository) Get(pkg string) (token apphost.AccessToken, err error) {
	if t.Dir == "" {
		t.Dir = dir.Token
	}
	if token, err = os.ReadJson[apphost.AccessToken](t.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}

func (t Repository) Set(pkg string, token apphost.AccessToken) (err error) {
	if t.Dir == "" {
		t.Dir = dir.Token
	}
	if err = os.WriteJson[apphost.AccessToken](token, t.Dir, pkg); err != nil {
		err = plog.Err(err)
	}
	return
}
