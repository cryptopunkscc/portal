package param

import "reflect"

type Values struct {
	Named      map[string]reflect.Value
	Positional []reflect.Value
}

func NewValues(args []any) (p *Values) {
	p = &Values{}
	p.Named = make(map[string]reflect.Value)
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

func (vs *Values) add(name string, v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Pointer:
		if !v.IsZero() { // TODO verify
			vs.add(name, v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			sf := v.Type().Field(i)
			tag := sf.Tag.Get("cli")
			if tag == "" && fv.Kind() != reflect.Struct {
				return
			}
			vs.add(tag, fv)
		}
	default:
		if name != "" {
			vs.Named[name] = v
		} else {
			vs.Positional = append(vs.Positional, v)
		}
	}
	return
}
