package goja

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/core/js/embed/common"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/dop251/goja"
)

type Backend struct {
	runtime *goja.Runtime
	core    bind.Core
	coreJs  string
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

func (b *Backend) Interrupt() {
	if b.runtime != nil {
		b.runtime.ClearInterrupt()
		b.core.Interrupt()
		b.runtime = nil
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
	b.Interrupt()
	b.runtime = goja.New()

	// bind core to goja runtime
	if err = Bind(b.runtime, b.core); err != nil {
		return plog.Err(err)
	}

	// inject core js adapter
	if _, err = b.runtime.RunString(b.coreJs); err != nil {
		return plog.Err(err)
	}

	// set args
	if err = b.runtime.Set("args", args); err != nil {
		return plog.Err(err)
	}

	// start js application
	if _, err = b.runtime.RunString(app); err != nil {
		return plog.Err(err)
	}
	return
}
