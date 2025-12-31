package source

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
)

func ObjectsCreateCommit(objects *astrald.ObjectsClient, object astral.Object) (objectID *astral.ObjectID, err error) {
	writer, err := objects.Create(nil, "", 0)
	if err != nil {
		return
	}
	if _, err = astral.DefaultBlueprints.Canonical().Write(writer, object); err != nil {
		return
	}
	return writer.Commit()
}
