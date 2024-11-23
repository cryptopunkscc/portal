package caller

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"reflect"
	"testing"
)

func TestStruct_Call(t *testing.T) {
	tests := []struct {
		name       string
		function   any
		unmarshal  Unmarshal
		data       []byte
		wantResult []any
		wantErr    bool
		assert     func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:       "1",
			assert:     assert.Equal,
			wantResult: []any{1, "foo", true},
			unmarshal: func(bytes []byte, args []any) error {
				reflect.ValueOf(args[0]).Elem().FieldByName("I").SetInt(1)
				reflect.ValueOf(args[1]).Elem().SetInt(1)
				reflect.ValueOf(args[2]).Elem().SetString("foo")
				reflect.ValueOf(args[3]).Elem().SetBool(true)
				return nil
			},
			function: func(t *testing.T, o testOptions, i int, s string, b bool, n any) (int, string, bool) {
				require.NotNil(t, t)
				require.Nil(t, n)
				return i, s, b
			},
		},
		{
			name: "2",
			function: func() <-chan int {
				c := make(chan int, 10)
				for i := 0; i < 10; i++ {
					c <- i
				}
				close(c)
				return c
			},
			assert: func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
				args := actual.([]any)
				assert.Len(t, args, 1)
				actual = args[0]
				if !assert.IsType(t, make(<-chan any), actual) {
					return false
				}
				c := actual.(<-chan any)
				i := 0
				for a := range c {
					if !assert.Equal(t, i, a) {
						return false
					}
					i++
				}
				return true
			},
		},
		{
			name:       "3",
			assert:     assert.Equal,
			wantResult: []any(nil),
			unmarshal: func(bytes []byte, args []any) error {
				return nil
			},
			function: func(o *testOptions) {
				log.Println(o)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.function)
			gotResult, err := c.Defaults(t).Unmarshalers(tt.unmarshal).Call(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.assert(t, tt.wantResult, gotResult)
		})
	}
}

type testOptions struct {
	I int    `cli:"i"`
	B bool   `cli:"b"`
	S string `cli:"s"`
}
