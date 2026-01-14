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
		objType := ""
		objType, _, err = astral.CanonicalTypeDecoder(reader)
		switch {
		case err != nil:
			return
		case objType != obj.ObjectType():
			return fmt.Errorf("object type mismatch: got %q, want %q", objType, obj.ObjectType())
		}
		_, err = obj.ReadFrom(reader)
	case io.ReaderFrom:
		_, err = obj.ReadFrom(reader)
	default:
		err = fmt.Errorf("%T is not io.ReaderFrom", v)
	}
	return
}
