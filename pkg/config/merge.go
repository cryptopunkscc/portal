package config

import "reflect"

func Merge[T any](cfg *T, others ...*T) {
	for _, other := range others {
		merge[T](cfg, other)
	}
}

// Merge recursively updates zero-value fields in cfg using values from defaultCfg.
func merge[T any](cfg, defaultCfg *T) {
	// Ensure we are dealing with pointers to structures.
	cfgVal := reflect.ValueOf(cfg).Elem()
	defVal := reflect.ValueOf(defaultCfg).Elem()

	mergeValue(cfgVal, defVal)
}

func mergeValue(cfgVal, defVal reflect.Value) {
	// Ensure both values are structures.
	if cfgVal.Kind() != reflect.Struct || defVal.Kind() != reflect.Struct {
		panic("config: cannot merge non-struct values")
	}

	numFields := cfgVal.NumField()
	for i := 0; i < numFields; i++ {
		field := cfgVal.Field(i)
		defField := defVal.Field(i)

		// If the field is a structure, perform recursive merging.
		if field.Kind() == reflect.Struct {
			mergeValue(field, defField)
			continue
		}

		// If the field is not set (is zero value), copy the value from the default.
		if isZeroOfUnderlyingType(field.Interface()) {
			// Check if the field can be set.
			if field.CanSet() {
				field.Set(defField)
			}
		}
	}
}

// isZeroOfUnderlyingType checks whether the given value is its zero value.
func isZeroOfUnderlyingType(x any) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
