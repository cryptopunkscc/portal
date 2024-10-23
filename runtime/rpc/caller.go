package rpc

import (
	"math"
	"reflect"
)

type Caller struct {
	name    string
	env     []any
	decoder argsDecoders
	f       reflect.Value
	args    []reflect.Value
}

func NewCaller(name string) (c *Caller) {
	c = &Caller{name: name}
	c.Decoder(NewJsonArgsDecoder(), NewClirArgsDecoder())
	return
}

func (c *Caller) With(env ...any) *Caller {
	cc := *c
	cc.env = append(c.env, env...)
	return &cc
}

func (c *Caller) Func(function any) *Caller {
	if c.f = reflect.ValueOf(function); c.f.Kind() != reflect.Func {
		panic("argument must be a function")
	}
	return c
}

func (c *Caller) Decoder(decoders ...ArgsDecoder) *Caller {
	return c.Decoders(decoders)
}

func (c *Caller) Decoders(decoders []ArgsDecoder) *Caller {
	c.decoder.Append(decoders)
	return c
}

func (c *Caller) Call(args ByteScannerReader) (out []any, err error) {
	values, err := c.call(args)
	if err != nil {
		return
	}
	out = formatOut(values)
	return
}

func (c *Caller) call(args ByteScannerReader) (out []reflect.Value, err error) {
	values, err := c.decodeIn(args)
	if err != nil {
		return
	}
	values = c.f.Call(values)
	err = handleError(values)
	if err != nil {
		return
	}
	out, err = c.runNested(values, args)
	return
}

func (c *Caller) decodeIn(args ByteScannerReader) (values []reflect.Value, err error) {
	injected := 0
	t := c.f.Type()

	// Inject dependencies
	var initial []reflect.Value
	for _, a := range c.env {
		initial = append(initial, reflect.ValueOf(a))
	}
	for i := 0; i < t.NumIn() && len(initial) > 0; i++ {
		for len(initial) > 0 && !initial[0].Type().AssignableTo(t.In(i)) {
			initial = initial[1:]
			continue
		}
		if len(initial) > 0 {
			values = append(values, initial[0])
			initial = initial[1:]
			injected++
		}
	}

	// prepare parameters for decoding
	var decoded []any
	for i := len(values); i < t.NumIn(); i++ {
		at := t.In(i)
		av := reflect.New(at)
		values = append(values, av.Elem())
		decoded = append(decoded, av.Interface())
	}

	variadicStartsAt := math.MaxInt

	// unfold varargs if needed
	if t.IsVariadic() {
		lastIndex := len(values) - 1
		decoded = decoded[:len(decoded)-1]
		last := values[lastIndex]
		values = values[:lastIndex]
		variadicStartsAt = lastIndex
		typ := last.Type().Elem()
		buffer := 20
		for i := 0; i < buffer; i++ {
			av := reflect.New(typ)
			values = append(values, av.Elem())
			decoded = append(decoded, av.Interface())
		}
	}

	// decode args
	if len(decoded) > 0 {
		if err = c.decoder.Decode(args, decoded); err != nil {
			return
		}
	}

	// trim empty varargs
	for i := variadicStartsAt; i < len(values); i++ {
		if values[i].IsZero() {
			values = values[0:i]
			break
		}
	}
	return
}

func handleError(values []reflect.Value) (err error) {
	if len(values) == 0 {
		return
	}
	last := values[len(values)-1]
	if last.Type().Implements(errorInterface) {
		if i := last.Interface(); i != nil {
			err, _ = i.(error)
		}
	}
	return
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func (c *Caller) runNested(values []reflect.Value, args ByteScannerReader) (r []reflect.Value, err error) {
	for _, value := range values {
		if value.Kind() == reflect.Func {
			e := *(&c)
			e.f = value
			var rr []reflect.Value
			if rr, err = e.call(args); err != nil {
				return
			}
			r = append(r, rr...)
			continue
		}
		r = append(r, value)
	}
	return
}

func formatOut(values []reflect.Value) (result []any) {
	// filter out error for values
	for _, value := range values {
		result = append(result, value.Interface())
	}

	// trim nil values
	for n := len(result) - 1; n > 0 && result[n] == nil; n-- {
		result = result[0:n]
	}
	return
}
