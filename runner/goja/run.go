package goja

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
	"reflect"
)

type Runner struct {
	bindings target.NewApi
	prefix   []string
}

func NewRunner(newApi target.NewApi, prefix ...string) target.Run[target.AppBackend] {
	return Runner{bindings: newApi, prefix: prefix}.Run
}

func (r Runner) Run(ctx context.Context, app target.AppBackend) (err error) {
	log.Println("Attach backend", reflect.TypeOf(app), app.Path(), app.Type())
	if err = NewBackend(r.bindings(ctx, app)).RunFs(app.Files()); err != nil {
		return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
	}
	<-ctx.Done()
	return
}
