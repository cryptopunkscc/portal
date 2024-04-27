package goja

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/common"
	"github.com/dop251/goja"
	"io/fs"
	"os"
	"path"
)

type Backend struct {
	vm        *goja.Runtime
	appHost   apphost.Flat
	appHostJs string
}

func NewBackend(ctx context.Context) *Backend {
	return &Backend{
		appHost:   apphost.WithTimeout(ctx),
		appHostJs: binding.CommonJsString,
	}
}

func (b *Backend) Run(app string) (err error) {
	if fs.ValidPath(app) {
		return b.RunPath(app)
	} else {
		return b.RunSource(app)
	}
}

func (b *Backend) RunPath(app string) (err error) {
	stat, err := os.Stat(app)
	if err != nil {
		return err
	}
	var src []byte
	if stat.IsDir() {
		app = path.Join(app, "main.js")
	}
	src, err = os.ReadFile(app)

	return b.RunSource(string(src))
}

func (b *Backend) RunFs(appFs fs.FS) (err error) {
	var src []byte
	if src, err = fs.ReadFile(appFs, "main.js"); err != nil {
		return err
	}
	return b.RunSource(string(src))
}

func (b *Backend) RunSource(app string) (err error) {
	if b.vm != nil {
		b.vm.ClearInterrupt()
		b.appHost.Interrupt()
	}
	b.vm = goja.New()

	if err = Bind(b.vm, b.appHost); err != nil {
		return
	}

	// inject apphost client js lib
	if _, err = b.vm.RunString(b.appHostJs); err != nil {
		return
	}

	// start js application backend
	if _, err = b.vm.RunString(app); err != nil {
		return
	}
	return
}
