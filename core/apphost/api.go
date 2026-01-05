package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/astrald"
)

func (a *Adapter) Dir() *astrald.DirClient {
	return astrald.NewDirClient(a.TargetID, a.Client)
}

func (a *Adapter) Objects() *astrald.ObjectsClient {
	return astrald.NewObjectsClient(a.TargetID, a.Client)
}
