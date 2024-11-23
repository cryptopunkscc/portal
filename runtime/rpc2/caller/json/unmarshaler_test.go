package json

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshaler_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		actual   []any
		expected []any
		wantErr  bool
	}{
		{
			data:     `[1, "a", true]`,
			expected: []any{float64(1), "a", true},
			actual:   []any{0, "", false},
		},
		{
			data:     `{"foo": "bar"}`,
			expected: []any{map[string]interface{}{"foo": "bar"}},
			actual:   []any{map[string]interface{}{}},
		},
		{
			data:     `[{"foo": "bar"}]`,
			expected: []any{map[string]interface{}{"foo": "bar"}},
			actual:   []any{map[string]interface{}{}},
		},
		{
			data:     `[{"foo": "bar"}]`,
			expected: []any{&Options{Foo: "bar"}},
			actual:   []any{&Options{}},
		},
	}
	u := Unmarshaler{}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
			if err := u.Unmarshal([]byte(tt.data), tt.actual); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.ElementsMatch(t, tt.expected, tt.actual)
		})
	}
}

type Options struct {
	Foo string `json:"foo"`
}
