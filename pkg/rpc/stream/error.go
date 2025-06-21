package stream

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Failure struct {
	Error error
}

// MarshalJSON returns m as the JSON encoding of m.
func (s *Failure) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"error": s.Error.Error(),
	})
}

// UnmarshalJSON sets *m to a copy of data.
func (s *Failure) UnmarshalJSON(data []byte) (err error) {
	m := map[string]string{}
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	ss, ok := m["error"]
	if !ok {
		return plog.Errorf("cannot umarshal json %s to %T", data, s)
	}
	if s == nil {
		*s = Failure{}
	}
	s.Error = errors.New(ss)
	return
}
