package target

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestAny(t *testing.T) {
	var counter = 0
	r := func(source Source) (Source, error) {
		counter++
		return source, nil
	}
	expected := testSource{}
	resolve := Any[Source](r, r, r)
	actual, _ := resolve(expected)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 1, counter)
}

func TestAny_nil(t *testing.T) {
	err := errors.New("test error")
	var counter = 0
	r := func(source Source) (Source, error) {
		counter++
		return source, err
	}
	resolve := Any[Source](r, r, r)
	actual, _ := resolve(testSource{})
	assert.Equal(t, nil, actual)
	assert.Equal(t, 3, counter)
}

type testSource struct{}

func (t testSource) Abs() (v string)                  { return }
func (t testSource) Path() (v string)                 { return }
func (t testSource) Files() (v fs.FS)                 { return }
func (t testSource) IsDir() (v bool)                  { return }
func (t testSource) Sub(string) (v Source, err error) { return }
