package golang

import (
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveProject(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedGo)
	bundle, err := ResolveProject(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedGoManifest, bundle.Manifest())
	assert.Equal(t, test.EmbedGoBuild, bundle.Build())
	assert.Equal(t, test.EmbedGoManifest, bundle.Dist().Manifest())
}

func TestResolveBundle(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedGoBundle)
	bundle, err := exec.ResolveBundle(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedGoManifest, bundle.Manifest())
}

func TestResolveDist(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedGoDist)
	bundle, err := exec.ResolveDist(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedGoManifest, bundle.Manifest())
}
