package basic

import (
	"context"
	"github.com/cryptopunkscc/portal/resolve/js"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"testing"
)

func TestGojaBackend(t *testing.T) {
	file, err := source.File("js")
	if err != nil {
		t.Fatal("can't resolve:", err)
	}
	dist, err := js.ResolveDist(file)
	if err != nil {
		return
	}
	runtime := bind.BackendRuntime()
	backend := goja.NewRunner(runtime)
	ctx := context.Background()
	err = backend.Run(ctx, dist, "foo", "bar")
	if err != nil {
		t.Error(err)
	}
}
