package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/out/common"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	frontend "github.com/cryptopunkscc/go-astral-js/pkg/webview"
	"github.com/webview/webview"
	"log"
	"os"
	"path"
)

func main() {
	file := os.Args[1]

	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	title := path.Base(file)
	src := string(srcBytes)

	ctx, cancel := context.WithCancel(context.Background())
	go exec.OnShutdown(cancel)

	w := frontend.New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(title)

	// inject apphost js client lib
	w.Init(binding.CommonJsString)

	// set app source code
	w.SetHtml(src)

	// bind apphost adapter to js env
	w.BindApphost(apphost.NewAdapter(ctx, nil))

	// start js application frontend
	w.Run()
}
