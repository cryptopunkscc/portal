package objects

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Dir struct {
	Path string
}

func (p Dir) Write(obj astral.Object) (objectID *astral.ObjectID, err error) {
	defer plog.TraceErr(&err)

	buf := bytes.NewBuffer(nil)
	if err = WriteCanonical(buf, obj); err != nil {
		return
	}
	if objectID, err = astral.ResolveObjectID(obj); err != nil {
		return
	}

	n := filepath.Join(p.Path, objectID.String())
	err = os.WriteFile(n, buf.Bytes(), 0644)
	return
}
