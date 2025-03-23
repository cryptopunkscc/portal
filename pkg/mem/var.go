package mem

import (
	"fmt"
	"reflect"
)

type String Val[string]
type MutableString Var[string]

type Val[T any] interface {
	Get() T
	Require() T
	IsZero() bool
	String() string
}

type Var[T any] interface {
	Val[T]
	Set(T)
}

type variable[T any] struct{ value T }

func (v *variable[T]) String() string { return fmt.Sprintf("%v", v.value) }

var _ Val[any] = NewVar[any](nil)

func NewVar[T any](value T) Var[T] {
	return &variable[T]{value: value}
}

func (v *variable[T]) Require() T {
	if v.IsZero() {
		panic("value required but not set")
	}
	return v.value
}

func (v *variable[T]) Get() T {
	return v.value
}

func (v *variable[T]) Set(value T) {
	v.value = value
}

func (v *variable[T]) IsZero() bool {
	return reflect.ValueOf(v.value).IsZero()
}
