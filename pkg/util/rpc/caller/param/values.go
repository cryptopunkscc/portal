package param

import (
	"reflect"
	"strings"
)

type Values struct {
	Type       string
	Named      map[string]reflect.Value
	Positional []reflect.Value
}

func NewValues(typ string, args []any) (p *Values) {
	p = &Values{}
	p.Type = typ
	p.Named = make(map[string]reflect.Value)
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		p.add("", v)
	}
	return
}

func (vs *Values) add(name string, v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Pointer:
		e := v.Elem()
		if e.Kind() == reflect.Pointer && e.IsZero() {
			n := reflect.New(e.Type().Elem())
			e.Set(n)
		}
		vs.add(name, e)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			sf := v.Type().Field(i)
			tag := sf.Tag.Get(vs.Type)
			if tag == "" && fv.Kind() != reflect.Struct {
				return
			}
			vs.add(tag, fv)
		}
	default:
		if name != "" {
			for _, n := range strings.Split(name, " ") {
				vs.Named[n] = v
			}
		} else {
			vs.Positional = append(vs.Positional, v)
		}
	}
}
