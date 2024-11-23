package clir

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
			data:     "1 2.2 a true",
			expected: []any{1, 2.2, "a", true},
			actual:   []any{0, 0.0, "", false},
		},
		{
			data:     `-i 1 -b true -s foo`,
			expected: []any{&TestArg{I: 1, B: true, TestArg2: TestArg2{S: "foo"}}},
			actual:   []any{&TestArg{}},
		},
		{
			data:     `-s foo -b true -i 1`,
			expected: []any{TestArg{I: 1, B: true, TestArg2: TestArg2{S: "foo"}}},
			actual:   []any{TestArg{}},
		},
		{
			data:     `1 true`,
			expected: []any{TestArgPos{I: 1, B: true}},
			actual:   []any{TestArgPos{}},
		},
		{
			data:     `1 true -s foo`,
			expected: []any{TestArgPos{I: 1, B: true, S: "foo"}},
			actual:   []any{TestArgPos{}},
		},
		{
			data:     `-s "foo" 1 true`,
			expected: []any{TestArgPos{I: 1, B: true, S: "foo"}},
			actual:   []any{TestArgPos{}},
		},
		{
			data:     `-s 'foo' 1 true`,
			expected: []any{TestArg2{S: "foo"}, 1, true},
			actual:   []any{TestArg2{}, 0, false},
		},
		{
			data:     `-s foo 1 true`,
			expected: []any{&TestArg2{S: "foo"}, 1, true},
			actual:   []any{&TestArg2{}, 0, false},
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

type TestArg struct {
	I int  `name:"i"`
	B bool `name:"b"`
	TestArg2
}

type TestArgPos struct {
	I int    `pos:"1"`
	B bool   `pos:"2"`
	S string `name:"s"`
}

type TestArg2 struct {
	S string `name:"s"`
}
