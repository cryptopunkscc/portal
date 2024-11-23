package caller

import (
	"github.com/stretchr/testify/require"
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
	}{
		{
			wantResult: []any{1, "foo", true},
			unmarshal: func(bytes []byte, args []any) error {
				reflect.ValueOf(args[0]).Elem().SetInt(1)
				reflect.ValueOf(args[1]).Elem().SetString("foo")
				reflect.ValueOf(args[2]).Elem().SetBool(true)
				return nil
			},
			function: func(t *testing.T, i int, s string, b bool, n any) (int, string, bool) {
				require.NotNil(t, t)
				require.Nil(t, n)
				return i, s, b
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
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Call() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
