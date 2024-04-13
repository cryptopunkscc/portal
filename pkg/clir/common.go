package clir

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/prod"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

func cliApplication(bindings runner.Bindings) func(f *FlagsPath) (err error) {
	return func(f *FlagsPath) (err error) {
		return prod.Run(bindings, f.Path)
	}
}
