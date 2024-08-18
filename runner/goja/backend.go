package goja

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/js/embed/common"
	"github.com/dop251/goja"
	"io/fs"
)

type Backend struct {
	vm        *goja.Runtime
	apphost   bind.Runtime
	apphostJs string
}

func NewBackend(apphost bind.Runtime) *Backend {
	if any(apphost) == nil {
		panic("apphost nil")
	}
	return &Backend{
		apphost:   apphost,
		apphostJs: common.JsString,
	}
}

func (b *Backend) RunFs(files fs.FS) (err error) {
	var src []byte
	if src, err = fs.ReadFile(files, "main.js"); err != nil {
		return plog.Err(err)
	}
	return b.RunSource(string(src))
}

func (b *Backend) RunSource(app string) (err error) {
	if b.vm != nil {
		b.vm.ClearInterrupt()
		b.apphost.Interrupt()
	}
	b.vm = goja.New()

	if err = Bind(b.vm, b.apphost); err != nil {
		return plog.Err(err)
	}

	// inject apphost client js lib
	if _, err = b.vm.RunString(b.apphostJs); err != nil {
		return plog.Err(err)
	}

	// start js application backend
	if _, err = b.vm.RunString(app); err != nil {
		return plog.Err(err)
	}
	return
}
