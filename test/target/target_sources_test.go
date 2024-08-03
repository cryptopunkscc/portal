package test

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test__sources_FromPath__test_assets(t *testing.T) {
	file, err := source.File("test_data")
	if err != nil {
		t.Fatal(err)
	}
	targets := target.List(sources.Resolver[target.Portal_](), file)

	for _, s := range targets {
		PrintTarget(s)
	}

	assert.LessOrEqual(t, 6, len(targets))
}

func Test__sources_FromFS__js_PortalLibFS(t *testing.T) {
	targets := target.List(sources.Resolver[target.Portal_](), source.Embed(js.PortalLibFS))

	for _, s := range targets {
		PrintTarget(s)
	}
}
