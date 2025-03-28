package astrald

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_Map(t *testing.T) {
	out := Config{}.Map()
	for n := range out {
		println(n)
	}
	assert.NotNil(t, out)
}
