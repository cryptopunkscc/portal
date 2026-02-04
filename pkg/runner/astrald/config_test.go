package astrald

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Map(t *testing.T) {
	out := Config{}.Map()
	for n := range out {
		println(n)
	}
	assert.NotNil(t, out)
}
