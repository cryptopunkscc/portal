package cli

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
	//var params = []any{opts, b, i, s, rest}
	var data = `-s "a b c" -i 10 -b true 1 's t r' lorem ipsum`

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

func TestUnmarshaler_Unmarshal_2(t *testing.T) {
	var opts *testOptions
	var params = []any{&opts}
	var data = `-s "a b c"`

	err := Unmarshal([]byte(data), params)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, opts, &testOptions{
		nestedOptions: nestedOptions{S: "a b c"},
	})
}

func TestUnmarshaler_Unmarshal_3(t *testing.T) {
	var opts *testOptions
	var params = []any{&opts}

	err := Unmarshal([]byte{}, params)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, opts, &testOptions{
		nestedOptions: nestedOptions{},
	})
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
