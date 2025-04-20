package query

import (
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		arg     any
		want    string
		wantErr bool
	}{
		{
			name: "named",
			arg: struct {
				I int
				S string
				B bool
			}{1, "s", true},
			want:    "b=true&i=1&s=s",
			wantErr: false,
		},
		{
			name: "map",
			arg: rpc.Opt{
				"b": true,
				"i": 1,
				"s": "s",
			},
			want:    "b=true&i=1&s=s",
			wantErr: false,
		},
		{
			name: "varargs",
			arg: []any{struct {
				I int
				S string
				B bool
			}{1, "s", true}, 1, "s", true},
			want: "_=1&_=s&_=true&b=true&i=1&s=s",
		},
		{
			name: "array",
			arg: []any{
				struct {
					A []int
				}{A: []int{1, 2, 3}},
			},
			want: "a=1&a=2&a=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := Marshal(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := string(bytes)
			if got != tt.want {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
