package apphost

import "github.com/cryptopunkscc/astrald/mod/dir/client"

func (a *Adapter) Dir() *dir.Client {
	return dir.New(a.TargetID, a.Client)
}
