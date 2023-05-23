package main

import (
	"astraljs"
	wails2 "astraljs/wails"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
)

func main() {
	if err := wails.Run(&options.App{
		Title:  "wails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Handler: wails2.NewFileLoader(astraljs.TryResolveWebApp()),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			&Adapter{*astraljs.NewAppHostFlatAdapter()},
		},
	}); err != nil {
		log.Fatalln(err)
	}
}

type Adapter struct{ astraljs.AppHostFlatAdapter }
