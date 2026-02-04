package fs2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanWriteToDir(t *testing.T) {
	assert.True(t, CanWriteToDir("./"))
	assert.False(t, CanWriteToDir("/bin"))
}
