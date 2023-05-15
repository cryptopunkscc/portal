package main

import (
	"astraljs"
	goja2 "astraljs/goja"
	"context"
	"github.com/dop251/goja"
	"log"
)

func main() {
	app := astraljs.ResolveWebApp()

	vm := goja.New()

	err := goja2.Bind(vm, astraljs.NewAppHostFlatAdapter())
	if err != nil {
		log.Fatal(err)
	}

	// inject apphost client js lib
	_, err = vm.RunString(astraljs.AppHostJsClient())
	if err != nil {
		log.Fatal(err)
	}

	// start js application backend
	_, err = vm.RunString(app.Source)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	<-ctx.Done()
}
