package object

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

func Unmarshal(b []byte, v any) (err error) {
	defer plog.TraceErr(&err)
	var r io.Reader = bytes.NewReader(b)
	switch t := v.(type) {
	case astral.Object:
		var typ string
		if typ, r, err = astral.OpenCanonical(r); err != nil {
			return
		}
		if typ != t.ObjectType() {
			return fmt.Errorf("object type mismatch: got %q, want %q", typ, t.ObjectType())
		}
		_, err = t.ReadFrom(r)
	case io.ReaderFrom:
		_, err = t.ReadFrom(r)
	default:
		err = fmt.Errorf("%T is not io.ReaderFrom", v)
	}
	return
}
