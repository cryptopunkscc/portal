package clir

import (
	"errors"
	"fmt"
	"github.com/leaanthony/clir"
	"reflect"
	"slices"
	"strings"
)

type Unmarshaler struct{}

func (u Unmarshaler) Unmarshal(data []byte, args []any) (err error) {
	fields := strings.Fields(string(data))

	// drop characters [", ']
	for i, field := range fields {
		if len(field) < 2 {
			continue
		}
		v := field[0]
		if v == '"' || v == '\'' {
			field = field[1:]
		}
		v = field[len(field)-1]
		if v == '"' || v == '\'' {
			field = field[:len(field)-1]
		}
		fields[i] = field
	}

	flags := false
	c := clir.NewCli("", "", "").Action(func() error { return nil })

	slices.Reverse(args)
	for i, arg := range args {
		value := reflect.ValueOf(arg)

		if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
			// add pointer to struct
			c.AddFlags(arg)
			flags = true
			continue
		}
		if value.Kind() == reflect.Struct {
			// add struct
			value = reflect.New(value.Type())
			c.AddFlags(value.Interface())
			args[i] = value.Elem()
			flags = true
			continue
		}

		if len(fields) == 0 {
			break
		}
		//s := fields[0]
		//fields = fields[1:]
		s := fields[len(fields)-1]
		fields = fields[:len(fields)-1]

		if value.Kind() == reflect.Ptr {
			// load pointer to primitive
			_, err = fmt.Sscan(s, arg)
		} else {
			// load primitive
			aa := reflect.New(value.Type())
			_, err = fmt.Sscan(s, aa.Interface())
			args[i] = aa.Elem().Interface()
		}

		if err == nil {
			continue
		}

		return errors.New("invalid arg type")
	}
	if flags {
		// load clir flags
		err = c.Run(fields...)
		if err != nil {
			return
		}

		// fix structs
		for i, arg := range args {
			switch v := arg.(type) {
			case reflect.Value:
				args[i] = v.Interface()
			}
		}
	}
	slices.Reverse(args)
	return
}
