package plog

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"slices"
	"testing"
)

func TestErr(t *testing.T) {
	expect := debug.Stack()
	actual := Err(errors.New("error")).stack

	slices.Reverse(expect)
	slices.Reverse(actual)

	same := 0
	for ; same < min(len(expect), len(actual)) && expect[same] == actual[same]; same++ {
	}

	assert.Equal(t, 201, same)
}
