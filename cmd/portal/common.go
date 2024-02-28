package main

import (
	astraljs "github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/prod"
	"io"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

func cliApplication(f *FlagsPath) (err error) {
	return prod.Run(f.Path, newAdapter)
}

type Adapter struct{ astraljs.FlatAdapter }

func newAdapter() io.Closer {
	return &Adapter{FlatAdapter: *astraljs.NewFlatAdapter()}
}
