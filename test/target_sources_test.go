package test

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	js "github.com/cryptopunkscc/go-astral-js/target/js/embed"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test__sources_FromPath__test_assets(t *testing.T) {
	assets := target.Abs("test_data")
	targets := sources.FromPath[target.Portal](assets)

	for _, s := range targets {
		PrintTarget(s)
	}

	assert.Equal(t, 6, len(targets))
}

func Test__sources_FromFS__js_PortalLibFS(t *testing.T) {
	targets := sources.FromFS[target.Source](js.PortalLibFS)

	for _, s := range targets {
		PrintTarget(s)
	}
}
