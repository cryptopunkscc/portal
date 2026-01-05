package goja

import (
	"context"
	"fmt"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/goja"
	source2 "github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/js"
)

type Runner interface {
	Run(ctx context.Context, args ...string) error
}

type BundleRunner struct {
	AppRunner
	Bundle js.JsBundle
}

func NewBundleRunner(core bind.Core) *BundleRunner {
	return &BundleRunner{AppRunner: AppRunner{Core: core}}
}

func (r *BundleRunner) ReadSrc(src source2.Source) (err error) {
	if err = r.Bundle.ReadSrc(src); err != nil {
		r.App = r.Bundle.App
		r.Func = r.Run
	}
	return
}

type AppRunner struct {
	js.App
	Core    bind.Core
	backend *goja.Backend
	args    []string
}

func NewAppRunner(core bind.Core) *AppRunner {
	return &AppRunner{Core: core}
}

func (r *AppRunner) ReadSrc(src source2.Source) (err error) {
	if err = r.App.ReadSrc(src); err != nil {
		r.Func = r.Run
	}
	return
}

func (r *AppRunner) Reload() (err error) {
	return r.backend.RunFs(r.FS(), r.args...)
}

func (r *AppRunner) Run(ctx context.Context, args ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", r.App.Metadata.Package, r.Path)
	r.args = args
	r.backend = goja.NewBackend(r.Core)
	if err = r.Reload(); err != nil {
		return
	}
	<-ctx.Done()
	r.backend.Interrupt()
	if r.Core.Code() > 0 {
		err = fmt.Errorf("exit %d", r.Core.Code())
	}
	return
}
