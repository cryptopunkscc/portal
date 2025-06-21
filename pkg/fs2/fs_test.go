package fs2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanWriteToDir(t *testing.T) {
	assert.True(t, CanWriteToDir("./"))
	assert.False(t, CanWriteToDir("/bin"))
}
