package wails

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"reflect"
)

type Runner struct {
	bindings target.New
	prefix   []string
}

func NewRunner(bindings target.New, prefix ...string) target.Run[target.AppFrontend] {
	return Runner{bindings: bindings, prefix: prefix}.Run
}

func (r Runner) Run(_ context.Context, app target.AppFrontend) (err error) {
	log.Println("Attach frontend", reflect.TypeOf(app), app.Path(), app.Type())
	opt := wails.AppOptions(r.bindings(target.TypeFrontend, r.prefix...))
	if err = wails.Run(app, opt); err != nil {
		return fmt.Errorf("dev.Run: %v", err)
	}
	return
}
