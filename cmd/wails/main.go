package main

import (
	"astraljs"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <app>\n", os.Args[0])
		os.Exit(0)
	}

	var app, err = NewApp(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := wails.Run(&options.App{
		Title:  app.Title,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Handler: app,
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
