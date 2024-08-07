package js

import (
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveDist(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedJsRollupDist)
	dist, err := ResolveDist(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedJsRollupManifest, dist.Manifest())
}

func TestResolveProject(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedJsRollup)
	project, err := ResolveProject(src)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedJsRollupManifest, project.Manifest())
	assert.Equal(t, test.EmbedJsRollupBuild, project.Build())
	assert.Equal(t, test.EmbedJsRollupManifest, project.Dist().Manifest())
}

func TestResolveBundle(t *testing.T) {
	defer test.Clean()
	src := test.Copy(test.EmbedJsRollupBundle)
	bundle, err := ResolveBundle(src)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedJsRollupManifest, bundle.Manifest())
}
