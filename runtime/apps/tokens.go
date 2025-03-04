package apps

import (
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/pkg/os"
)

type Tokens struct {
	Dir string
}

func (t Tokens) GetToken(pkg string) (token apphost.AccessToken, err error) {
	if t.Dir == "" {
		t.Dir = Dir
	}
	return os.ReadJson[apphost.AccessToken](t.Dir, pkg)
}

func (t Tokens) SetToken(pkg string, token apphost.AccessToken) (err error) {
	if t.Dir == "" {
		t.Dir = Dir
	}
	return os.WriteJson[apphost.AccessToken](token, t.Dir, pkg)
}
