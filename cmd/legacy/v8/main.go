package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/common"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/v8"
	"log"
	"os"
	"rogchap.com/v8go"
)

func main() {
	file := os.Args[1]

	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	src := string(srcBytes)

	iso := v8go.NewIsolate()
	defer iso.Dispose()

	// prepare context
	ctx, cancel := context.WithCancel(context.Background())
	exec.OnShutdown(cancel)

	// bind apphost adapter to js env
	global, err := v8.Bind(iso, apphost.NewAdapter(ctx))
	if err != nil {
		log.Fatal(err)
	}

	// create v8 context with app host bindings
	v8Ctx := v8go.NewContext(iso, global)
	defer v8Ctx.Close()

	// inject apphost client js lib
	_, err = v8Ctx.RunScript(binding.CommonJsString, "apphost")
	if err != nil {
		log.Fatal(err)
	}

	// start js application backend
	_, err = v8Ctx.RunScript(src, file)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
