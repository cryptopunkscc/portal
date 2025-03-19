package query

import (
	"errors"
	"github.com/cryptopunkscc/portal/core/rpc/caller/param"
	"net/url"
	"reflect"
	"strconv"
)

func Unmarshal(data []byte, params []any) (err error) {
	q, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	p := param.NewValues("query", params)
	err = setData(p, q)
	return
}

func setData(p *param.Values, v url.Values) (err error) {
	if positional, ok := v["_"]; ok {
		l := len(p.Positional)
		for i, s := range positional {
			if i == l {
				break
			}
			if err = setValue(p.Positional[i], s); err != nil {
				return
			}
		}
	}

	for key, values := range v {
		if value, ok := p.Named[key]; ok {
			err = setValue(value, values...)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func setValue(val reflect.Value, strs ...string) error {
	switch val.Kind() {
	case reflect.String:
		if len(strs) > 0 {
			val.SetString(strs[0])
		}
	case reflect.Int:
		if len(strs) > 0 {
			i, err := strconv.Atoi(strs[0])
			if err != nil {
				return err
			}
			val.SetInt(int64(i))
		}
	case reflect.Int8:
		if len(strs) > 0 {
			i, err := strconv.ParseInt(strs[0], 10, 8)
			if err != nil {
				return err
			}
			val.SetInt(i)
		}
	case reflect.Int16:
		if len(strs) > 0 {
			i, err := strconv.ParseInt(strs[0], 10, 16)
			if err != nil {
				return err
			}
			val.SetInt(i)
		}
	case reflect.Int32:
		if len(strs) > 0 {
			i, err := strconv.ParseInt(strs[0], 10, 32)
			if err != nil {
				return err
			}
			val.SetInt(i)
		}
	case reflect.Int64:
		if len(strs) > 0 {
			i, err := strconv.ParseInt(strs[0], 10, 64)
			if err != nil {
				return err
			}
			val.SetInt(i)
		}
	case reflect.Uint:
		if len(strs) > 0 {
			i, err := strconv.ParseUint(strs[0], 10, 0)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}
	case reflect.Uint8:
		if len(strs) > 0 {
			i, err := strconv.ParseUint(strs[0], 10, 8)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}
	case reflect.Uint16:
		if len(strs) > 0 {
			i, err := strconv.ParseUint(strs[0], 10, 16)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}
	case reflect.Uint32:
		if len(strs) > 0 {
			i, err := strconv.ParseUint(strs[0], 10, 32)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}
	case reflect.Uint64:
		if len(strs) > 0 {
			i, err := strconv.ParseUint(strs[0], 10, 64)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}
	case reflect.Float32:
		if len(strs) > 0 {
			f, err := strconv.ParseFloat(strs[0], 32)
			if err != nil {
				return err
			}
			val.SetFloat(f)
		}
	case reflect.Float64:
		if len(strs) > 0 {
			f, err := strconv.ParseFloat(strs[0], 64)
			if err != nil {
				return err
			}
			val.SetFloat(f)
		}
	case reflect.Bool:
		if len(strs) > 0 {
			if strs[0] == "" {
				val.SetBool(true)
				return nil
			}
			b, err := strconv.ParseBool(strs[0])
			if err != nil {
				return err
			}
			val.SetBool(b)
		}
	case reflect.Slice:
		newSlice := reflect.MakeSlice(val.Type(), len(strs), len(strs))
		for i := 0; i < len(strs); i++ {
			if err := setValue(newSlice.Index(i), strs[i]); err != nil {
				return err
			}
		}
		val.Set(newSlice)
	case reflect.Array:
		for i := 0; i < val.Len() && i < len(strs); i++ {
			if err := setValue(val.Index(i), strs[i]); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported type")
	}
	return nil
}
