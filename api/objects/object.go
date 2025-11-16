package objects

import (
	"fmt"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func ReadCanonical(reader io.Reader, obj astral.Object) (err error) {
	defer plog.TraceErr(&err)
	objType, _, err := astral.ReadCanonicalType(reader)
	switch {
	case err != nil:
		return
	case objType.String() != obj.ObjectType():
		return fmt.Errorf("invalid object type: %s", objType)
	}
	_, err = obj.ReadFrom(reader)
	return
}

func WriteCanonical(writer io.Writer, obj astral.Object) (err error) {
	defer plog.TraceErr(&err)
	_, err = astral.DefaultBlueprints.Canonical().Write(writer, obj)
	return
}
