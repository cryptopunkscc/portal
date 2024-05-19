package goja

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"reflect"
)

type Runner struct {
	bindings target.New
	prefix   []string
}

func NewRunner(bindings target.New, prefix ...string) target.Run[target.AppBackend] {
	return Runner{bindings: bindings, prefix: prefix}.Run
}

func (r Runner) Run(ctx context.Context, app target.AppBackend) (err error) {
	log.Println("Attach backend", reflect.TypeOf(app), app.Path(), app.Type())
	if err = goja.NewBackend(r.bindings(target.TypeBackend, r.prefix...)).RunFs(app.Files()); err != nil {
		return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
	}
	<-ctx.Done()
	return
}
