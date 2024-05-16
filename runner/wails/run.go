package wails

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"reflect"
)

func Run(
	bindings runtime.New,
	app target.App,
	prefix ...string,
) (err error) {
	log.Println("Attach frontend", reflect.TypeOf(app), app.Path(), app.Type())
	opt := wails.AppOptions(bindings(target.Frontend, prefix...))
	if err = wails.Run(app, opt); err != nil {
		return fmt.Errorf("dev.Run: %v", err)
	}
	return
}
