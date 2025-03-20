package goja

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/core/js/embed/common"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/dop251/goja"
	"io/fs"
)

type Backend struct {
	vm     *goja.Runtime
	core   bind.Core
	coreJs string
}

func NewBackend(core bind.Core) *Backend {
	if any(core) == nil {
		panic("apphost nil")
	}
	return &Backend{
		core:   core,
		coreJs: common.JsString,
	}
}

func (b *Backend) RunFs(files fs.FS, args ...string) (err error) {
	var src []byte
	if src, err = fs.ReadFile(files, "main.js"); err != nil {
		return plog.Err(err)
	}
	return b.RunSource(string(src), args...)
}

func (b *Backend) RunSource(app string, args ...string) (err error) {
	if b.vm != nil {
		b.vm.ClearInterrupt()
		b.core.Interrupt()
	}
	b.vm = goja.New()

	if err = Bind(b.vm, b.core); err != nil {
		return plog.Err(err)
	}

	// inject apphost client js lib
	if _, err = b.vm.RunString(b.coreJs); err != nil {
		return plog.Err(err)
	}

	// set args
	if err = b.vm.Set("args", args); err != nil {
		return plog.Err(err)
	}

	// start js application backend
	if _, err = b.vm.RunString(app); err != nil {
		return plog.Err(err)
	}
	return
}
