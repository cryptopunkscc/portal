package object

import (
	"bytes"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Marshal(v any) (b []byte, err error) {
	defer plog.TraceErr(&err)
	buf := bytes.NewBuffer(b)
	switch t := v.(type) {
	case astral.Object:
		_, err = astral.Encode(buf, t, astral.Canonical())
	case io.WriterTo:
		_, err = t.WriteTo(buf)
	default:
		err = plog.Errorf("%T is not io.WriterTo", v)
	}
	b = buf.Bytes()
	return
}
