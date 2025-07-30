package cli

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshaler_Unmarshal(t *testing.T) {
	plog.Verbosity = 100
	var opts testOptions
	var b bool
	var i int
	var s string
	var rest string
	var params = []any{&opts, &b, &i, &s, &rest}
	//var data = `-s "a b c" -i 10 -bp true 1 's t r' lorem ipsum`
	var data = `-s "a b c" -i 10 -b -p true 1 's t r' lorem ipsum`

	err := Unmarshal([]byte(data), params)
	if err != nil {
		plog.Println(err)
		t.Fatal()
	}

	assert.Equal(t, testOptions{
		nestedOptions: nestedOptions{S: "a b c"},
		I:             10,
		B:             true,
		P:             true,
	}, opts)
	assert.Equal(t, true, b)
	assert.Equal(t, 1, i)
	assert.Equal(t, `s t r`, s)
	assert.Equal(t, `lorem ipsum`, rest)
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
	P bool `cli:"p"`
}

type nestedOptions struct {
	S string `cli:"s"`
	T *testOptions
}
