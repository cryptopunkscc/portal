package v8

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/target/js/embed/common"
	"log"
	"rogchap.com/v8go"
)

func Run(ctx context.Context, file, src string) (err error) {
	iso := v8go.NewIsolate()
	defer iso.Dispose()

	// bind apphost adapter to js env
	global, err := Bind(iso, apphost.NewAdapter(ctx, nil))
	if err != nil {
		log.Fatal(err)
	}

	// create v8 context with app host bindings
	v8Ctx := v8go.NewContext(iso, global)
	defer v8Ctx.Close()

	// inject apphost client js lib
	_, err = v8Ctx.RunScript(common.JsString, "apphost")
	if err != nil {
		log.Fatal(err)
	}

	// start js application backend
	_, err = v8Ctx.RunScript(src, file)
	if err != nil {
		log.Fatal(err)
	}
	return
}