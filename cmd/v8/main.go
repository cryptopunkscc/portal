package main

import (
	"astraljs"
	"astraljs/v8"
	"context"
	"log"
	"rogchap.com/v8go"
)

func main() {
	app := astraljs.ResolveWebApp()

	iso := v8go.NewIsolate()
	defer iso.Dispose()

	// bind apphost adapter to js env
	global, err := v8.Bind(iso, astraljs.NewAppHostFlatAdapter())
	if err != nil {
		log.Fatal(err)
	}

	// create v8 context with app host bindings
	v8Ctx := v8go.NewContext(iso, global)
	defer v8Ctx.Close()

	// inject apphost client js lib
	_, err = v8Ctx.RunScript(astraljs.AppHostJsClient(), "apphost")
	if err != nil {
		log.Fatal(err)
	}

	// start js application backend
	_, err = v8Ctx.RunScript(app.Source, app.Path)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	<-ctx.Done()
}
