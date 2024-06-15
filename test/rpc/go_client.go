package rpc

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClient struct {
	port string
}

func NewTestClient(port string) *TestClient {
	return &TestClient{port: port}
}

func (c TestClient) Run(t *testing.T) {
	t.Log("Starting test client")

	services := []string{
		"go",
		"js",
	}

	tests := []struct {
		name    string
		getConn func(string, *testing.T) rpc.Conn
		skip    bool
	}{
		{
			//skip: true,
			name: "request",
			getConn: func(srv string, t *testing.T) rpc.Conn {
				return rpc.NewRequest(id.Anyone, fmt.Sprintf(c.port, srv), "request")
			},
		},
		{
			//skip: true,
			name: "flow",
			getConn: func(srv string, t *testing.T) (conn rpc.Conn) {
				conn, err := rpc.QueryFlow(id.Anyone, fmt.Sprintf(c.port, srv), "flow")
				if err != nil {
					t.Skip(err)
				}
				return
			},
		},
	}

	for _, srv := range services {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if test.skip {
					t.SkipNow()
				}

				request := test.getConn(srv, t)
				request.Logger(plog.Type(request))

				t.Run("func1", func(t *testing.T) {
					t.Run("a", func(t *testing.T) {
						str, err := rpc.Query[string](request, "func1", "text", false)
						assert.Nil(t, err)
						assert.Equal(t, "text", str)
					})
					t.Run("b", func(t *testing.T) {
						str, err := rpc.Query[string](request, "func1", "text", true)
						assert.Equal(t, "", str)
						assert.Equal(t, errors.New("text"), err)
					})
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

				t.Run("func3", func(t *testing.T) {
					t.Run("a", func(t *testing.T) {
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

					t.Run("b", func(t *testing.T) {
						actual, err := rpc.Query[*TestStruct2](request, "func3", nil)
						if err != nil {
							t.Fatal(err)
						}
						assert.Zero(t, actual)
					})
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
			})
		}
	}
}
