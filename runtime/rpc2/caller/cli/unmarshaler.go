package cli

import (
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/param"
	"reflect"
	"strings"
)

type Unmarshaler struct{}

func (u Unmarshaler) Unmarshal(data []byte, args []any) (err error) {
	return Unmarshal(data, args)
}

func (u Unmarshaler) Score(data []byte) (score uint) {
	for _, r := range string(data) {
		switch r {
		case ' ', '-':
			score++
		}
	}
	return
}

func Unmarshal(data []byte, args []any) (err error) {
	p := param.NewValues(args)
	f := parseFields(data)
	err = set(p, 0, f)
	return
}

func parseFields(data []byte) (fields []string) {
	var inString *rune = nil
	fields = strings.FieldsFunc(string(data), func(r rune) bool {
		if inString == nil {
			if r == '"' || r == '\'' {
				inString = &r
				return true
			}
			return r == ' '
		}
		if *inString == r {
			inString = nil
			return true
		}
		return false
	})
	return
}

func set(p *param.Values, offset int, fields []string) (err error) {
	if offset == len(fields) {
		return
	}
	field := fields[offset]
	var value reflect.Value
	if field[0] == '-' {
		value = p.Named[field[1:]]
		offset++
		if value.Kind() == reflect.Bool {
			value.SetBool(true)
			return set(p, offset, fields)
		}
		field = fields[offset]
	} else {
		if len(p.Positional) == 0 {
			return
		}
		value = p.Positional[0]
		p.Positional = p.Positional[1:]

		// special case to consume rest fields as string
		if len(p.Positional) == 0 && value.Kind() == reflect.String {
			value.SetString(strings.Join(fields[offset:], " "))
			return
		}
	}
	if value.Kind() == reflect.String {
		value.SetString(field)
	} else {
		if value.CanAddr() {
			value = value.Addr()
		}
		_, err = fmt.Sscan(field, value.Interface())
		if err != nil {
			return
		}
	}
	offset++
	return set(p, offset, fields)
}
