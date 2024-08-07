package html

import (
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveDist(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedSvelteDist)
	dist, err := ResolveDist(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedSvelteManifest, dist.Manifest())
}

func TestResolveProject(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedSvelte)
	project, err := ResolveProject(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedSvelteManifest, project.Manifest())
	assert.Equal(t, test.EmbedSvelteBuild, project.Build())
	assert.Equal(t, test.EmbedSvelteManifest, project.Dist().Manifest())
}

func TestResolveBundle(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedSvelteBundle)
	bundle, err := ResolveBundle(src)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedSvelteManifest, bundle.Manifest())
}
