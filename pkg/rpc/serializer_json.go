package rpc

import (
	"encoding/json"
	"io"
)

func JsonCodecs(rw io.ReadWriter) (e Encoder, d Decoder, m Marshal, u Unmarshal) {
	e = json.NewEncoder(rw)
	d = json.NewDecoder(rw)
	m = json.Marshal
	u = json.Unmarshal
	return
}

func (a *raw) UnmarshalJSON(b []byte) error {
	a.bytes = b
	return nil
}

func (a *raw) MarshalJSON() ([]byte, error) {
	return a.bytes, nil
}

type Failure struct {
	Error string `json:"error"`
}
