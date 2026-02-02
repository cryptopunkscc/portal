package player

import (
	"encoding/json"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
)

type Status struct {
	ObjectID *astral.ObjectID
	Position astral.Duration
	Length   astral.Duration
}

func (Status) ObjectType() string { return "mediaplayer.status" }

func (s Status) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Struct(s).WriteTo(w)
}

func (s *Status) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Struct(s).ReadFrom(r)
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func init() {
	_ = astral.Add(&Status{})
}
