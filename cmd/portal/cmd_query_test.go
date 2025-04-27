package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "options",
			query: "foo bar -a 1 -b 2",
			want:  "foo.bar?a=1&b=2",
		},
		{
			name:  "varargs",
			query: "foo bar - a 1 -b 2",
			want:  "foo.bar?_=a&_=1&_=-b&_=2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, parseQuery(tt.query), "parseQuery(%v)", tt.query)
		})
	}
}
