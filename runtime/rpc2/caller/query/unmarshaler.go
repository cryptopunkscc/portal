package query

import (
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller/param"
	"net/url"
	"reflect"
	"strings"
)

type Unmarshaler struct{}

func (u Unmarshaler) Unmarshal(data []byte, params []any) (err error) {
	return Unmarshal(data, params)
}

func (u Unmarshaler) Score(data []byte) (score uint) {
	for _, r := range string(data) {
		switch r {
		case '&', '=':
			score++
		}
	}
	return
}

func Unmarshal(data []byte, params []any) (err error) {
	e, err := parseArgs(string(data))
	if err != nil {
		return
	}
	p := param.NewValues(params)
	err = set(p, 0, e)
	return
}

type arg struct {
	key    string
	values []string
}

func parseArgs(query string) (args []arg, err error) {
	indexes := make(map[string]int)
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		if i, b := indexes[key]; !b {
			indexes[key] = len(args)
			args = append(args, arg{key, []string{value}})
		} else {
			args[i].values = append(args[i].values, value)
		}
	}
	return
}

func set(p *param.Values, offset int, elements []arg) (err error) {
	if offset == len(elements) {
		return
	}

	entry := elements[offset]
	key := entry.key
	args := entry.values
	offset++
	var field string
	var value reflect.Value
	if len(args) > 0 && args[0] != "" {
		ok := false
		value, ok = p.Named[key]
		if !ok {
			return fmt.Errorf("unrecognized option '%s'", key)
		}
		field = args[0]
	} else {
		if v, ok := p.Named[key]; ok {
			if v.Kind() == reflect.Bool {
				v.SetBool(true)
				return set(p, offset, elements)
			}
		}

		if len(p.Positional) == 0 {
			return
		}
		field = key
		value = p.Positional[0]
		p.Positional = p.Positional[1:]
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
	return set(p, offset, elements)
}
