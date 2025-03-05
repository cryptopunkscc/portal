package tokens

import (
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/os"
)

type Repository struct {
	Dir string
}

func (t Repository) Get(pkg string) (token apphost.AccessToken, err error) {
	if t.Dir == "" {
		t.Dir = Dir
	}
	return os.ReadJson[apphost.AccessToken](t.Dir, pkg)
}

func (t Repository) Set(pkg string, token apphost.AccessToken) (err error) {
	if t.Dir == "" {
		t.Dir = Dir
	}
	return os.WriteJson[apphost.AccessToken](token, t.Dir, pkg)
}
