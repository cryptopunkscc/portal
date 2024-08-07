package dist

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Any(target.Source) (any, error) {
	return 1, nil
}

func TestResolve_EmbedJs(t *testing.T) {
	a, err := Resolver(Any).Resolve(test.EmbedJs)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedJsManifest, a.Manifest())
	assert.Equal(t, 1, a.Target())
}

func TestResolve_EmbedHtml(t *testing.T) {
	a, err := Resolver(Any).Resolve(test.EmbedHtml)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedHtmlManifest, a.Manifest())
	assert.Equal(t, 1, a.Target())
}

func TestResolve_EmbedJsRollup(t *testing.T) {
	_, err := Resolver(Any).Resolve(test.EmbedJsRollup)
	if err != nil {
		assert.ErrorIs(t, err, target.ErrNotTarget)
	}
}

func TestResolve_EmbedSvelte(t *testing.T) {
	_, err := Resolver(Any).Resolve(test.EmbedSvelte)
	if err != nil {
		assert.ErrorIs(t, err, target.ErrNotTarget)
	}
}

func TestResolve_List(t *testing.T) {
	expected := map[string]*target.Manifest{
		"data.zip/go/dist":        test.EmbedGoManifest,
		"data.zip/html":           test.EmbedHtmlManifest,
		"data.zip/js":             test.EmbedJsManifest,
		"data.zip/js-rollup/dist": test.EmbedJsRollupManifest,
		"data.zip/sh":             test.EmbedShManifest,
		"data.zip/svelte/dist":    test.EmbedSvelteManifest,
	}
	for _, dist := range Resolver(Any).List(test.EmbedRoot) {
		s, ok := expected[dist.Abs()]
		if !ok {
			t.Error("unexpected target", dist.Manifest().Package)
		}
		assert.Equal(t, s, dist.Manifest())
		delete(expected, dist.Abs())
	}
	assert.Equal(t, map[string]*target.Manifest{}, expected)
}
