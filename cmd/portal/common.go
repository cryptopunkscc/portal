package main

import (
	"context"
	astraljs "github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/backend/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/frontend/wails"
	"github.com/wailsapp/wails/v2/pkg/options"
	"sync"
)

func AppOptions() *options.App {
	adapter := &Adapter{*astraljs.NewFlatAdapter()}
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []interface{}{adapter},
		OnDomReady: func(ctx context.Context) {
			adapter.Interrupt()
		},
	}
}

type Adapter struct{ astraljs.FlatAdapter }

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

type FlagsApp struct {
	FlagsPath
	Front bool `name:"f"`
	Back  bool `name:"b"`
}

func (f *FlagsApp) Setup() {
	if !f.Front && !f.Back {
		f.Front = true
		f.Back = true
	}
}

func cliApplication(f *FlagsApp) (err error) {
	f.Setup()
	wait := sync.WaitGroup{}
	if f.Back {
		wait.Add(1)
		if err = goja.NewBackend().Run(f.Path); err != nil {
			return
		}
	}
	if f.Front {
		wait.Add(1)
		return wails.Run(f.Path, AppOptions())
	}
	wait.Wait()
	return
}
