package goja

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"reflect"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	app target.App,
	prefix ...string,
) (err error) {
	log.Println("Attach backend", reflect.TypeOf(app), app.Path(), app.Type())
	if err = goja.NewBackend(bindings(target.Backend, prefix...)).RunFs(app.Files()); err != nil {
		return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
	}
	<-ctx.Done()
	return
}
