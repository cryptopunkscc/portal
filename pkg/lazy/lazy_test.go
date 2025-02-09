package lazy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	f := V(func() *foo { return &foo{} })

	f().a = "foo"

	assert.Equal(t, "foo", f().a)
}

type foo struct {
	a string
}
