package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/template"
	"github.com/stretchr/testify/assert"
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
