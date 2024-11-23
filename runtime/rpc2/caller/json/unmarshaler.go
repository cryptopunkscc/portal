package json

import (
	"encoding/json"
	"errors"
)

type Unmarshaler struct{}

func (u Unmarshaler) Score(data []byte) (score uint) {
	for _, r := range string(data) {
		switch r {
		case '[', ']', '{', '}', '"', ':':
			score++
		}
	}
	return
}

func (u Unmarshaler) Unmarshal(data []byte, args []any) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}
	if data[0] != '{' && data[0] != '[' {
		return errors.New("not a JSON object")
	}
	// unmarshal struct payload to as first arg
	if len(args) == 1 && data[0] == '{' {
		return json.Unmarshal(data, &args[0])
	}
	// perform default unmarshal
	return json.Unmarshal(data, &args)
}
