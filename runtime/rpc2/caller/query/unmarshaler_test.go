package query

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {

	tests := []struct {
		name    string
		data    string
		params  []any
		expect  []any
		wantErr bool
	}{
		{
			name: "test1",
			data: `s=a b c&int=10&b=true`,
			params: func() []any {
				var opts testOptions
				return []any{&opts}
			}(),
			expect: []any{
				testOptions{
					nestedOptions: nestedOptions{S: "a b c"},
					I:             10,
					B:             true,
				},
			},
		},
		{
			name: "test2",
			data: `_=1&_=2&_=3`,
			params: func() []any {
				var i1 int
				var i2 int
				var i3 int
				return []any{&i1, &i2, &i3}
			}(),
			expect: []any{1, 2, 3},
		},
		{
			name: "test3",
			data: `s=a b c&int=10&b=&_=true&_=1&_=s t r&_=lorem ipsum`,
			params: func() []any {
				var opts testOptions
				var b bool
				var i int
				var s string
				var rest string
				return []any{&opts, &b, &i, &s, &rest}
			}(),
			expect: []any{
				testOptions{
					nestedOptions: nestedOptions{S: "a b c"},
					I:             10,
					B:             true,
				},
				true,
				1,
				"s t r",
				"lorem ipsum",
			},
		},
		{
			name: "test4",
			data: `_=1&_=2`,
			params: func() []any {
				var i1 int
				var i2 int
				var i3 int
				return []any{&i1, &i2, &i3}
			}(),
			expect: []any{1, 2, 0},
		},
		{
			name: "test5",
			data: `_=1&_=2&_=3`,
			params: func() []any {
				var i1 int
				var i2 int
				return []any{&i1, &i2}
			}(),
			expect: []any{1, 2},
		},
		{
			name: "test6",
			data: "_=.%2Ffoo%2Fbar%2Fbaz%2F",
			params: func() []any {
				var path string
				return []any{&path}
			}(),
			expect: []any{"./foo/bar/baz/"},
		},
		{
			name:   "test7",
			data:   `a=1&a=2&a=3`,
			params: func() []any { return []any{&testSlice{}} }(),
			expect: []any{testSlice{[]int{1, 2, 3}}},
		},
		{
			name:   "test8",
			data:   `a=1&a=2&a=3`,
			params: func() []any { return []any{&testArray{}} }(),
			expect: []any{testArray{[3]int{1, 2, 3}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal([]byte(tt.data), tt.params); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				for i, param := range tt.params {
					tt.params[i] = reflect.ValueOf(param).Elem().Interface()
				}
				assert.Equal(t, tt.expect, tt.params)
			}
		})
	}
}

type testOptions struct {
	nestedOptions
	I int  `query:"int i"`
	B bool `query:"b"`
}

type nestedOptions struct {
	S string `query:"s"`
	T *testOptions
}

type testSlice struct {
	A []int `query:"a"`
}

type testArray struct {
	A [3]int `query:"a"`
}
