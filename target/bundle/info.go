package bundle

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/manifest"
)

type Info struct {
	Manifest manifest.App
	Release
	ReleaseID *astral.ObjectID
	Host      *astral.Identity
}
