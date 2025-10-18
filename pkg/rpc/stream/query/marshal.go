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

	switch objType.Kind() {
	case reflect.Array, reflect.Slice:
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

	case reflect.Map:
		if objValue.Len() == 0 {
			return ErrorEmptyValue
		}
		for _, key := range objValue.MapKeys() {
			mapValue := objValue.MapIndex(key)
			// skip nil pointers in map values
			if mapValue.Kind() == reflect.Ptr && mapValue.IsNil() {
				continue
			}
			// skip zero values
			if mapValue.IsZero() {
				continue
			}

			if err = marshal(mapValue.Interface(), to, key.String()); err != nil {
				return
			}
		}

	case reflect.Struct:
		if objValue.NumField() == 0 {
			return ErrorEmptyValue
		}
		for i := 0; i < objType.NumField(); i++ {
			field := objType.Field(i)
			fieldValue := objValue.Field(i)

			// skip nil pointers
			if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				continue
			}

			// Skip unexported fields
			if !fieldValue.CanInterface() {
				continue
			}

			includeEmpty := strings.Contains(field.Tag.Get("include"), "empty")

			// skip zero values
			if !includeEmpty && fieldValue.IsZero() {
				continue
			}

			// use the "query" tag if available, otherwise fallback to the field name
			fieldTag := field.Tag.Get("query")
			if fieldTag == "" || fieldTag == "-" {
				fieldTag = strings.ToLower(field.Name)
			}

			// handle collections
			if fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice {
				if err = marshal(fieldValue.Interface(), to, fieldTag); err != nil {
					return
				}
				continue
			}

			// format single value
			value := fmt.Sprintf("%v", fieldValue.Interface())
			if includeEmpty || value != "" {
				to.Add(fieldTag, value)
			}
		}

	default:
		// for basic types, add value with a fallback key
		to.Add(tag, fmt.Sprintf("%v", from))
	}

	return
}
