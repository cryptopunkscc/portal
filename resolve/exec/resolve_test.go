package exec

import (
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveDist(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedSh)
	a, err := ResolveDist(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedShManifest, a.Manifest())
}
