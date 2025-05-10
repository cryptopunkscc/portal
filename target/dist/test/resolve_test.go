package test

import (
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve(t *testing.T) {
	expected := manifest.App{}
	err := all.Unmarshalers.Unmarshal(PortalYaml, &expected)
	test.AssertErr(t, err)

	s, err := source.Embed(DistFS).Sub("test_dist")
	test.AssertErr(t, err)

	p, err := dist.Resolve_(s)
	test.AssertErr(t, err)
	assert.Equal(t, expected, *p.Manifest())
}
