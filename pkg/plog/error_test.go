package plog

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"slices"
	"testing"
)

func TestErr_stack(t *testing.T) {
	expect := debug.Stack()
	actual := Err(errors.New("error")).stack

	slices.Reverse(expect)
	slices.Reverse(actual)

	same := 0
	for ; same < min(len(expect), len(actual)) && expect[same] == actual[same]; same++ {
	}

	assert.Equal(t, 201, same)
}

func TestErr_Error_1(t *testing.T) {
	actual := Err(errors.New("error"), "foo")

	assert.Equal(t, "foo: error", actual.Error())
}

func TestErr_Error_2(t *testing.T) {
	actual := Err(nil, "foo")

	assert.Equal(t, "foo", actual.Error())
}
