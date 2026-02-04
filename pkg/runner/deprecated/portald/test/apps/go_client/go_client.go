package go_client

import (
	"testing"
)

type TestClient struct {
	port     string
	services []string
}

func NewTestClient(port string, services ...string) *TestClient {
	if len(services) == 0 {
		panic("must provide at least one service")
	}
	return &TestClient{port: port, services: services}
}

func (c TestClient) Run(t *testing.T) {
	//log.Println("Starting test client")
	//
	//tests := []struct {
	//	name    string
	//	getConn func(string, *testing.T) apphost.Conn
	//	skip    bool
	//}{
	//	{
	//		//skip: true,
	//		name: "request",
	//		getConn: func(srv string, t *testing.T) apphost.Conn {
	//			//query := apphost2.FormatPort(c.port, srv, "request")
	//			//return apphost.Default.Rpc().Request("localnode", query)
	//			panic("implement me")
	//		},
	//	},
	//	{
	//		//skip: true,
	//		name: "flow",
	//		getConn: func(srv string, t *testing.T) (conn apphost.Conn) {
	//			//query := apphost2.FormatPort(c.port, srv, "flow")
	//			//conn, err := apphost.Default.Rpc().Conn("localnode", query)
	//			//if err != nil {
	//			//	t.Skip(err)
	//			//}
	//			//return
	//			panic("implement me")
	//
	//		},
	//	},
	//}
	//
	//for _, srv := range c.services {
	//	t.Run(srv, func(t *testing.T) {
	//		for _, test := range tests {
	//			t.Run(test.name, func(t *testing.T) {
	//				if test.skip {
	//					t.SkipNow()
	//				}
	//
	//				request := test.getConn(srv, t)
	//				request.Logger(plog.Type(request))
	//
	//				t.Run("func1", func(t *testing.T) {
	//					t.Run("a", func(t *testing.T) {
	//						str, err := rpc.Query[string](request, "func1", "text", false)
	//						assert.Nil(t, err)
	//						assert.Equal(t, "text", str)
	//					})
	//					t.Run("b", func(t *testing.T) {
	//						str, err := rpc.Query[string](request, "func1", "text", true)
	//						assert.Equal(t, "", str)
	//						assert.Equal(t, errors.New("RPC: text"), err)
	//					})
	//				})
	//
	//				t.Run("func2", func(t *testing.T) {
	//					expected := []any{true, float64(1), 99.99, "text"}
	//					actual, err := rpc.Query[[]any](request, "func2", expected...)
	//					if err != nil {
	//						t.Fatal(err)
	//					}
	//					assert.Equal(t, expected, actual)
	//				})
	//
	//				t.Run("func3", func(t *testing.T) {
	//					t.Run("a", func(t *testing.T) {
	//						expected := rpc2.TestStruct1{
	//							B: false,
	//							I: 1,
	//							F: 0,
	//							S: "",
	//						}
	//
	//						actual, err := rpc.Query[rpc2.TestStruct1](request, "func3", expected)
	//						if err != nil {
	//							t.Fatal(err)
	//						}
	//						assert.Equal(t, expected, actual)
	//					})
	//
	//					t.Run("b", func(t *testing.T) {
	//						actual, err := rpc.Query[*rpc2.TestStruct1](request, "func3", []byte{})
	//						assert.Zero(t, actual)
	//						assert.Equal(t, query.ErrorEmptyValue, err)
	//					})
	//				})
	//
	//				t.Run("func4", func(t *testing.T) {
	//					arg := []any{true, 1, 99.99, "text"}
	//					expected := rpc2.TestStruct1{
	//						B: true,
	//						I: 1,
	//						F: 99.99,
	//						S: "text",
	//					}
	//					actual, err := rpc.Query[rpc2.TestStruct1](request, "func4", arg...)
	//					if err != nil {
	//						t.Fatal(err)
	//					}
	//					assert.Equal(t, expected, actual)
	//				})
	//			})
	//		}
	//	})
	//}
}
