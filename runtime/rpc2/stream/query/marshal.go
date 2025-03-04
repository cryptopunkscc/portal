package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func Marshal(v any) (bytes []byte, err error) {
	vals := url.Values{}
	if err = marshal(v, vals, "_"); err == nil {
		bytes = []byte(vals.Encode())
	}
	return
}

var ErrorEmptyValue = fmt.Errorf("cannot marshal empty value as uri query")

func marshal(from any, to url.Values, tag string) (err error) {
	if from == nil {
		return ErrorEmptyValue
	}
	objType := reflect.TypeOf(from)
	objValue := reflect.ValueOf(from)

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}

	if objType.Kind() == reflect.Array || objType.Kind() == reflect.Slice {
		if objValue.Len() == 0 {
			return ErrorEmptyValue
		}

		for i := 0; i < objValue.Len(); i++ {
			element := objValue.Index(i)
			err = marshal(element.Interface(), to, tag)
			if err != nil {
				return
			}

		}
	} else if objType.Kind() == reflect.Struct {
		if objValue.NumField() == 0 {
			return ErrorEmptyValue
		}
		for i := 0; i < objType.NumField(); i++ {
			field := objType.Field(i)
			fieldValue := objValue.Field(i)
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				continue
			}
			if fieldValue.IsZero() {
				continue
			}

			// use the `json` tag if available, otherwise fallback to the field name
			tag = field.Tag.Get("query")
			if tag == "" || tag == "-" {
				tag = field.Name
				tag = strings.ToLower(tag)
			}

			// Skip unexported fields
			if !fieldValue.CanInterface() {
				continue
			}

			// handle collections
			if fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
				if err = marshal(fieldValue.Interface(), to, tag); err != nil {
					return
				}
				continue
			}

			// format single value
			value := fmt.Sprintf("%v", fieldValue.Interface())
			if value != "" {
				to.Add(tag, value)
			}
		}
	} else {
		to.Add(tag, fmt.Sprintf("%v", from))
	}

	return
}
