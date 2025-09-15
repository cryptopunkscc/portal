package app

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/source"
	"log"
	"os"
	"path/filepath"
)

type Installer struct {
	Dir          string
	Prepare      func(target.App_) error
	Repositories target.Repositories
	Resolvers    []target.Resolve[target.Source]
}

var _ target.Run[target.App_] = Installer{}.Run

func (i Installer) Dispatcher() target.Dispatcher {
	return target.Dispatcher{
		Runner:   target.RunSeq,
		Provider: i.runnableProvider(),
	}
}

func (i Installer) runnableProvider() target.Provider[target.Runnable] {
	return target.Provider[target.Runnable]{
		Repository: i.Repositories,
		Resolve:    target.Any[target.Runnable](i.runner().Try),
		Filter:     canRun,
	}
}

func (i Installer) runner() *target.SourceRunner[target.App_] {
	return &target.SourceRunner[target.App_]{
		Runner: i,
		Resolve: func(src target.Source) (out target.App_, err error) {
			resolvers := append(i.Resolvers, Resolve_.Try)
			return target.Any[target.App_](resolvers...).Resolve(src)
		},
	}
}

func canRun(app target.Runnable) bool {
	if e, ok := app.Source().(target.AppExec); ok {
		return e.Runtime().Target().Match()
	}
	return true
}

func (i Installer) Run(_ context.Context, src target.App_, _ ...string) (err error) {
	return i.Install(src)
}

func (i Installer) Install(src target.App_) (err error) {
	defer plog.TraceErr(&err)
	plog.Println("installing", src.Manifest().Package, src.Manifest().Runtime)
	if err = i.prepare(src); err != nil {
		return
	}
	p := i.dstPath(src)
	err = src.CopyTo(p)
	return
}

func (i Installer) prepare(src target.App_) (err error) {
	if i.Prepare != nil {
		return i.Prepare(src)
	}
	return
}

func (i Installer) dstPath(src target.App_) (out string) {
	return filepath.Join(i.Dir, src.Manifest().Package)
}

func (i Installer) Uninstall(id string) (err error) {
	defer plog.TraceErr(&err)
	src := source.Dir(i.Dir)
	for _, t := range Resolve_.List(src) {
		if t.Manifest().Match(id) {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			return os.RemoveAll(t.Abs())
		}
	}
	return fmt.Errorf("%s not found", id)
}
