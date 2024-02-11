package goja

import (
	astraljs "github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/dop251/goja"
	"io/fs"
	"log"
)

type Backend struct {
	vm      *goja.Runtime
	appHost *astraljs.FlatAdapter
}

func NewBackend() *Backend {
	return &Backend{
		appHost: astraljs.NewFlatAdapter(),
	}
}

func (b *Backend) Run(path string) (err error) {
	// identify app bundle type
	bundleType, err := assets.BundleType(path)
	if err != nil {
		return
	}

	bundleFs, err := assets.BundleFS(bundleType, path)
	if err != nil {
		return
	}

	bytes, err := fs.ReadFile(bundleFs, "service.js")
	if err != nil {
		return err
	}

	b.RunSource(string(bytes))
	return
}

func (b *Backend) RunSource(app string) {
	if b.vm != nil {
		b.vm.ClearInterrupt()
		b.appHost.Interrupt()
	}
	b.vm = goja.New()

	err := Bind(b.vm, b.appHost)
	if err != nil {
		log.Fatal(err)
	}

	// inject apphost client js lib
	_, err = b.vm.RunString(astraljs.JsBaseString())
	if err != nil {
		log.Fatal(err)
	}

	// start js application backend
	_, err = b.vm.RunString(app)
	if err != nil {
		log.Fatal(err)
	}
}
