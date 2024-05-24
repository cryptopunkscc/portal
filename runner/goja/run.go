package goja

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"reflect"
)

func NewRun(newApi target.NewApi, prefix ...string) target.Run[target.AppBackend] {
	return Runner{bindings: newApi, prefix: prefix}.Run
}

type Runner struct {
	bindings target.NewApi
	prefix   []string
}

func (r Runner) Run(ctx context.Context, app target.AppBackend) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("Attach backend", reflect.TypeOf(app), app.Path(), app.Type())
	if err = NewBackend(r.bindings(ctx, app)).RunFs(app.Files()); err != nil {
		return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
	}
	<-ctx.Done()
	return
}
