package apphost

import "github.com/cryptopunkscc/astrald/lib/astrald"

func (a *Adapter) Dir() *astrald.DirClient {
	return astrald.NewDirClient(a.TargetID, a.Client)
}
