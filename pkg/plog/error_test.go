package plog

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"slices"
	"testing"
)

func init() {
	Verbosity = all
}

func TestErr_stack(t *testing.T) {
	expect := debug.Stack()
	actual := Err(errors.New("error")).stack

	println("=======================================")
	println(string(expect))
	println("=======================================")
	println(string(actual))
	println("=======================================")

	slices.Reverse(expect)
	slices.Reverse(actual)

	same := 0
	for ; same < min(len(expect), len(actual)) && expect[same] == actual[same]; same++ {
	}
	println(fmt.Sprintln("same suffix bytes:", same))
	assert.Greater(t, same, 0)
}

func TestErr_Error_1(t *testing.T) {
	actual := Err(errors.New("error"), "foo")

	assert.Equal(t, "foo: error", actual.Error())
}

func TestErr_Error_2(t *testing.T) {
	actual := Err(nil, "foo")

	assert.Equal(t, "foo", actual.Error())
}

func TestTraceErr(t *testing.T) {
	err := errorWithTrace()
	assert.ErrorAs(t, err, &ErrStack{})
}

func TestPrintln(t *testing.T) {
	err := errorWithTrace()
	Println(err)
}

func errorWithTrace() (err error) {
	defer TraceErr(&err)
	return errors.New("error")
}
