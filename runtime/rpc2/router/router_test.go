package router

import (
	"github.com/cryptopunkscc/portal/runtime/rpc2/caller"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"github.com/cryptopunkscc/portal/runtime/rpc2/registry"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRouter_Query(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		route    string
		expected *cmd.Handler
	}{
		{
			name:     "load caller if query is correct",
			query:    "foo bar",
			route:    "foo",
			expected: &cmd.Handler{Func: func() {}},
		},
		{
			name:     "load noting if query isn't correct",
			query:    "foo bar",
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Base{
				Registry: registry.New[*cmd.Handler]().Add(tt.route, tt.expected),
			}.Query(tt.query)

			actual := router.Registry.Get()
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestRouter_Call(t *testing.T) {
	tests := []struct {
		name      string
		route     string
		args      string
		caller    *cmd.Handler
		unmarshal caller.Unmarshaler
		deps      []any
		expected  []any
	}{
		{
			name:     "function should be called",
			expected: []any{[]any{1, true}},
			caller:   &cmd.Handler{Func: func() (int, bool) { return 1, true }},
		},
		{
			name:   "unmarshaler should be called",
			args:   "foo",
			caller: &cmd.Handler{Func: func(_ string) {}},
			unmarshal: caller.Unmarshal(func(data []byte, args []any) error {
				assert.Equal(t, "foo", string(data))
				return nil
			}),
			expected: []any{stream.End},
		},
		{
			name: "dependencies should be passed",
			deps: []any{t},
			caller: &cmd.Handler{Func: func(tt *testing.T) {
				assert.NotNil(tt, t)
			}},
			expected: []any{stream.End},
		},
		{
			name: "nil error should be omitted",
			caller: &cmd.Handler{Func: func() error {
				return nil
			}},
			expected: []any{stream.End},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Base{
				Registry:     registry.New[*cmd.Handler]().Add(tt.route, tt.caller),
				Unmarshalers: []caller.Unmarshaler{tt.unmarshal},
				Dependencies: tt.deps,
				args:         tt.args,
			}.Call()
			var actual []any
			for a := range c {
				actual = append(actual, a)
			}

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
