package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type Unmarshaler struct{}

func (u Unmarshaler) Unmarshal(data []byte, args []any) (err error) {
	return Unmarshal(data, args)
}

func Unmarshal(data []byte, args []any) (err error) {
	p := newParams(args)
	f := parseFields(data)
	err = p.set(0, f)
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

type values struct {
	named      map[string]reflect.Value
	positional []reflect.Value
}

func newParams(args []any) (p *values) {
	p = &values{}
	p.named = make(map[string]reflect.Value)
	//va := reflect.ValueOf(args)
	//for i := 0; i < va.Len(); i++ {
	//	v := va.Index(i)
	//	p.add("", v.Elem())
	//
	//}
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

func (p *values) set(offset int, fields []string) (err error) {
	if offset == len(fields) {
		return
	}
	field := fields[offset]
	var value reflect.Value
	if field[0] == '-' {
		value = p.named[field[1:]]
		offset++
		if value.Kind() == reflect.Bool {
			value.SetBool(true)
			return p.set(offset, fields)
		}
		field = fields[offset]
	} else {
		if len(p.positional) == 0 {
			return
		}
		value = p.positional[0]
		p.positional = p.positional[1:]

		// special case to consume rest fields as string
		if len(p.positional) == 0 && value.Kind() == reflect.String {
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
	return p.set(offset, fields)
}
