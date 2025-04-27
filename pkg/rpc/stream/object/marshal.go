package object

import (
	"bytes"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

func Marshal(v any) (b []byte, err error) {
	defer plog.TraceErr(&err)
	buf := bytes.NewBuffer(b)
	switch t := v.(type) {
	case astral.Object:
		_, err = astral.WriteCanonical(buf, t)
	case io.WriterTo:
		_, err = t.WriteTo(buf)
	default:
		err = plog.Errorf("%T is not io.WriterTo", v)
	}
	b = buf.Bytes()
	return
}
