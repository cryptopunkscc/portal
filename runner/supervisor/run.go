package supervisor

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"sync"
)

func Runner[T target.Portal_](
	wait *sync.WaitGroup,
	processes *sig.Map[string, T],
	run target.Run[T],
) target.Run[T] {
	return runner[T]{
		wait:    wait,
		running: processes,
		run:     run,
	}.Run
}

type runner[T target.Portal_] struct {
	wait    *sync.WaitGroup
	running *sig.Map[string, T]
	run     target.Run[T]
}

func (r runner[T]) Run(ctx context.Context, portal T) (err error) {
	log := plog.Get(ctx)
	id := portal.Manifest().Package
	log.Println("setting", id)
	if _, ok := r.running.Set(id, portal); !ok {
		log.Println(r.running.Clone())
		log.Printf("%s already started ", id)
		return
	}
	r.wait.Add(1)
	log.Printf("start %T %s %s", portal, portal.Manifest().Package, portal.Abs())
	err = r.run(ctx, portal)
	log.Printf("exit %T %s %s", portal, portal.Manifest().Package, portal.Abs())
	r.running.Delete(id)
	r.wait.Done()
	return
}
