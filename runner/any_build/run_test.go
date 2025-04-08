package any_build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"log"
	"sync"
	"testing"
)

func Test_provider(t *testing.T) {
	plog.Verbosity = 100
	for i, runnable := range provider.Provide("../../apps") {
		log.Printf("%d. %T %s", i, runnable.Source(), runnable.Manifest().Package)
	}
}

func Test_dispatch_sync(t *testing.T) {
	t.SkipNow()
	plog.Verbosity = 100
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	d := target.Dispatcher{
		Provider: provider,
		Runner:   target.RunSeq,
	}
	err := d.Run(ctx, "../../apps", "pack", "clean")
	wg.Wait()
	if err != nil {
		plog.Println(err)
		t.Error(err)
	}
}

func Test_dispatch_async(t *testing.T) {
	plog.Verbosity = 100
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	d := target.Dispatcher{
		Provider: provider,
		Runner:   &target.AsyncRunner{WaitGroup: wg},
	}
	err := d.Run(ctx, "../../apps", "pack", "clean")
	wg.Wait()
	if err != nil {
		plog.Println(err)
		t.Error(err)
	}
}
