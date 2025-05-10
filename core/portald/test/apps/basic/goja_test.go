package basic

import (
	"context"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/source"
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
	core := bind.NewBackendCore
	runner := goja.NewRunner(core)
	ctx := context.Background()
	err = runner.Run(ctx, dist, "foo", "bar")
	if err != nil {
		t.Error(err)
	}
}
