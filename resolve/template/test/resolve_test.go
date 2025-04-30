package test

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve(t *testing.T) {
	d, err := source.Embed(templateFS).Sub("template")
	test.AssertErr(t, err)

	s, err := template.Resolve(d)
	if err != nil {
		return
	}

	assert.Equal(t, "name", s.Info().Name)
}
