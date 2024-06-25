package goja

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/js/embed/common"
	"github.com/dop251/goja"
	"io/fs"
	"os"
	"path"
)

type Backend struct {
	vm        *goja.Runtime
	apphost   target.Apphost
	apphostJs string
}

func NewBackend(apphost target.Apphost) *Backend {
	return &Backend{
		apphost:   apphost,
		apphostJs: common.JsString,
	}
}

func (b *Backend) Run(app string) (err error) {
	return b.RunPath(app)
}

func (b *Backend) RunPath(app string) (err error) {
	stat, err := os.Stat(app)
	if err != nil {
		return plog.Err(err)
	}
	var src []byte
	if stat.IsDir() {
		app = path.Join(app, "main.js")
	}
	src, err = os.ReadFile(app)
	if err != nil {
		return plog.Err(err)
	}

	return b.RunSource(string(src))
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
