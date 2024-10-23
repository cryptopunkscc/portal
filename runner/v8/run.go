package v8

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/runtime/js/embed/common"
	"log"
	"rogchap.com/v8go"
)

func Run(ctx context.Context, file, src string) (err error) {
	iso := v8go.NewIsolate()
	defer iso.Dispose()

	// bind apphost adapter to js env
	var runtime bind.Runtime
	//ah = apphost.NewFactory(nil).WithTimeout(ctx, "src") // FIXME
	global, err := Bind(iso, runtime)
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
