package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type Unmarshaler struct{}

func (u Unmarshaler) Unmarshal(data []byte, params []any) (err error) {
	return Unmarshal(data, params)
}

func Unmarshal(data []byte, params []any) (err error) {
	e, err := parseArgs(string(data))
	if err != nil {
		return
	}
	p := newValues(params)
	err = p.set(0, e)
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

type values struct {
	named      map[string]reflect.Value
	positional []reflect.Value
}

func newValues(args []any) (p *values) {
	p = &values{}
	p.named = make(map[string]reflect.Value)
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		p.add("", v)
	}
	return
}

func (p *values) add(name string, v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Pointer:
		if !v.IsZero() { // TODO verify
			p.add(name, v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			sf := v.Type().Field(i)
			tag := sf.Tag.Get("cli")
			if tag == "" && fv.Kind() != reflect.Struct {
				return
			}
			p.add(tag, fv)
		}
	default:
		if name != "" {
			p.named[name] = v
		} else {
			p.positional = append(p.positional, v)
		}
	}
	return
}

func (p *values) set(offset int, elements []arg) (err error) {
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
		value = p.named[key]
		field = args[0]
	} else {
		if v, ok := p.named[key]; ok {
			if v.Kind() == reflect.Bool {
				v.SetBool(true)
				return p.set(offset, elements)
			}
		}

		if len(p.positional) == 0 {
			return
		}
		field = key
		value = p.positional[0]
		p.positional = p.positional[1:]
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
	return p.set(offset, elements)
}
