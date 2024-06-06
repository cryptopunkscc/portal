package rpc

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClientTest struct {
	port string
}

func NewTestClientTest(port string) *TestClientTest {
	return &TestClientTest{port: port}
}

func (c TestClientTest) Run(t *testing.T) {
	t.Log("Starting test client")
	request := rpc.NewRequest(id.Anyone, c.port)

	t.Run("func1", func(t *testing.T) {
		if err := rpc.Call(request, "func1"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("func2", func(t *testing.T) {
		expected := []any{true, float64(1), 99.99, "text"}
		actual, err := rpc.Query[[]any](request, "func2", expected...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("func3a", func(t *testing.T) {
		expected := TestStruct2{TestStruct1{
			B: false,
			I: 1,
			F: 0,
			S: "",
		}}

		actual, err := rpc.Query[TestStruct2](request, "func3", expected)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("func3b", func(t *testing.T) {
		actual, err := rpc.Query[*TestStruct2](request, "func3", nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Zero(t, actual)
	})

	t.Run("func4", func(t *testing.T) {
		arg := []any{true, 1, 99.99, "text"}
		expected := TestStruct1{
			B: true,
			I: 1,
			F: 99.99,
			S: "text",
		}
		actual, err := rpc.Query[TestStruct1](request, "func4", arg...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(actual)
		assert.Equal(t, expected, actual)
	})

}
