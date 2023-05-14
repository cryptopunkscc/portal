package main

import (
	astral_js "astral-js"
	goja2 "astral-js/goja"
	"context"
	"github.com/dop251/goja"
	"log"
)

func main() {
	app := astral_js.ResolveWebApp()

	vm := goja.New()

	err := goja2.Bind(vm, astral_js.NewAppHostFlatAdapter())
	if err != nil {
		log.Fatal(err)
	}

	// inject apphost client js lib
	_, err = vm.RunString(astral_js.AppHostJsClient())
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
