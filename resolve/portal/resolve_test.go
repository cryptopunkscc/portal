package portal

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve_EmbedJs(t *testing.T) {
	a, err := Resolve[any](test.EmbedJs)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedJsManifest, a.Manifest())
}

func TestResolve_EmbedHtml(t *testing.T) {
	a, err := Resolve[any](test.EmbedHtml)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, test.EmbedHtmlManifest, a.Manifest())
}

func TestResolve_EmbedJsRollup(t *testing.T) {
	_, err := Resolve[any](test.EmbedJsRollup)
	if err != nil {
		assert.ErrorIs(t, err, target.ErrNotTarget)
	}
}

func TestResolve_EmbedSvelte(t *testing.T) {
	_, err := Resolve[any](test.EmbedSvelte)
	if err != nil {
		assert.ErrorIs(t, err, target.ErrNotTarget)
	}
}
