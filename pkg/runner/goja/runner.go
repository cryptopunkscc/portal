package goja

import (
	"context"
	"fmt"

	"github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/source/js"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

type Runner interface {
	Run(ctx *bind.Core, args ...string) error
}

type BundleRunner struct {
	AppRunner
	Bundle js.Bundle
}

func (r BundleRunner) New() source.Source {
	return &r
}

func (r *BundleRunner) ReadSrc(src source.Source) (err error) {
	if err = r.Bundle.ReadSrc(src); err == nil {
		r.App = r.Bundle.App
		r.Func = r.Run
	}
	return
}

type AppRunner struct {
	js.App
	Core    *bind.Core
	backend *Backend
	Args    []string
}

func (r AppRunner) New() source.Source {
	return &r
}

func (r *AppRunner) ReadSrc(src source.Source) (err error) {
	if err = r.App.ReadSrc(src); err != nil {
		r.Func = r.Run
	}
	return
}

func (r *AppRunner) Reload(ctx context.Context) (err error) {
	r.Core.Interrupt()
	return r.backend.RunFs(r.PathFS(), r.Args...)
}

func (r *AppRunner) Start(ctx *bind.Core, args ...string) (err error) {
	r.Core = ctx
	log := plog.Get(ctx).Type(r)
	log.Printf("run %T %s", r.App.Metadata.Package, r.Path)
	r.Args = args
	r.backend = NewBackend(r.Core)
	return r.Reload(ctx)
}
func (r *AppRunner) Run(ctx *bind.Core, args ...string) (err error) {
	if err = r.Start(ctx, args...); err != nil {
		return
	}
	<-ctx.Done()
	r.backend.Interrupt()
	if ctx.Code() > 0 {
		err = fmt.Errorf("exit %d", ctx.Code())
	}
	return
}
