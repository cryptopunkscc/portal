package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshaler_Unmarshal(t *testing.T) {
	var opts testOptions
	var b bool
	var i int
	var s string
	var rest string
	var params = []any{&opts, &b, &i, &s, &rest}
	var data = `s=a b c&int=10&b=true&true&1&s t r&lorem ipsum`

	err := Unmarshal([]byte(data), params)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, opts, testOptions{
		nestedOptions: nestedOptions{S: "a b c"},
		I:             10,
		B:             true,
	})
	assert.Equal(t, b, true)
	assert.Equal(t, i, 1)
	assert.Equal(t, s, `s t r`)
	assert.Equal(t, rest, `lorem ipsum`)
}

type testOptions struct {
	nestedOptions
	I int  `cli:"int i"`
	B bool `cli:"b"`
}

type nestedOptions struct {
	S string `cli:"s"`
	T *testOptions
}
