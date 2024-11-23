package caller

import (
	"fmt"
	"math"
	"reflect"
)

type Func struct {
	Names        []string
	function     reflect.Value
	defaults     []any
	unmarshalers []Unmarshaler
}

func New(function any, unmarshalers ...Unmarshaler) (c *Func) {
	c = new(Func)
	if c.function = reflect.ValueOf(function); c.function.Kind() != reflect.Func {
		panic("not a function")
	}
	c.unmarshalers = unmarshalers
	return
}

func (c *Func) Named(names ...string) *Func {
	c.Names = names
	return c
}

func (c *Func) Unmarshalers(unmarshalers ...Unmarshaler) *Func {
	cc := &(*c)
	if len(unmarshalers) > 0 {
		cc.unmarshalers = unmarshalers
	}
	return cc
}

func (c *Func) Defaults(defaults ...any) *Func {
	cc := &(*c)
	cc.defaults = append(cc.defaults, defaults...)
	return cc
}

func (c *Func) Call(data []byte) (result []any, err error) {
	values, err := c.invoke(data)
	if err != nil {
		return
	}
	result = extractResult(values)
	return
}

func (c *Func) invoke(data []byte) (out []reflect.Value, err error) {
	values, err := c.decodeArguments(data)
	if err != nil {
		return
	}
	values = c.function.Call(values)
	err = getError(values)
	if err != nil {
		return
	}
	out, err = c.runNested(values, data)
	return
}

func (c *Func) decodeArguments(data []byte) (values []reflect.Value, err error) {
	injected := 0
	t := c.function.Type()

	// Inject dependencies
	var initial []reflect.Value
	for _, a := range c.defaults {
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
		if err = c.unmarshal(data, decoded); err != nil {
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

func (c *Func) unmarshal(data []byte, args []any) error {
	if len(c.unmarshalers) == 0 {
		panic("no unmarshalers specified")
	}
	chosen := c.unmarshalers[0]
	score := chosen.Score(data)
	for i := 1; i < len(c.unmarshalers); i++ {
		n := c.unmarshalers[i]
		s := n.Score(data)
		if score < s {
			chosen = n
			score = s
		}
	}
	err := chosen.Unmarshal(data, args)
	if err == nil {
		return nil
	}
	return fmt.Errorf("cannot unmarshal data %s %v", data, err)
}

func getError(values []reflect.Value) (err error) {
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

func (c *Func) runNested(values []reflect.Value, data []byte) (r []reflect.Value, err error) {
	for _, value := range values {
		if value.Kind() == reflect.Func {
			e := *c
			e.function = value
			var rr []reflect.Value
			if rr, err = e.invoke(data); err != nil {
				return
			}
			r = append(r, rr...)
			continue
		}
		r = append(r, value)
	}
	return
}

func extractResult(values []reflect.Value) (result []any) {
	// filter out error for values
	for _, value := range values {
		var add any
		if value.Kind() == reflect.Chan {
			add = valueToChan(value)
		} else {
			add = value.Interface()
		}
		result = append(result, add)
	}

	// trim nil values
	for n := len(result) - 1; n > 0 && result[n] == nil; n-- {
		result = result[0:n]
	}
	return
}

func valueToChan(value reflect.Value) <-chan any {
	out := make(chan any)
	go func() {
		sel := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: value}}
		defer close(out)
		for {
			if _, v, b := reflect.Select(sel); b {
				out <- v.Interface()
			} else {
				return
			}
		}
	}()
	return out
}
