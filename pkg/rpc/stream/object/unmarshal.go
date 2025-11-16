package object

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Unmarshal(b []byte, v any) (err error) {
	defer plog.TraceErr(&err)
	var reader io.Reader = bytes.NewReader(b)
	switch obj := v.(type) {
	case astral.Object:
		var objType astral.ObjectType
		objType, _, err = astral.ReadCanonicalType(reader)
		switch {
		case err != nil:
			return
		case objType.String() != obj.ObjectType():
			return fmt.Errorf("object type mismatch: got %q, want %q", objType.String(), obj.ObjectType())
		}
		_, err = obj.ReadFrom(reader)
	case io.ReaderFrom:
		_, err = obj.ReadFrom(reader)
	default:
		err = fmt.Errorf("%T is not io.ReaderFrom", v)
	}
	return
}
