package go_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/flow"
	golang "github.com/cryptopunkscc/go-astral-js/pkg/go"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"time"
)

type Runner struct {
	watcher *golang.Watcher
	send    target.MsgSend
	build   Build
	run     target.Run[target.DistExec]

	dist   target.DistExec
	ctx    context.Context
	cancel context.CancelFunc
}

type Build func(context.Context, ...string) error

func NewRunner(
	build Build,
	run target.Run[target.DistExec],
	send target.MsgSend,
) *Runner {
	return &Runner{
		watcher: golang.NewWatcher(),
		run:     run,
		build:   build,
		send:    send,
	}
}

func (r *Runner) Reload() (err error) {
	if r.cancel != nil {
		r.cancel()
	}
	ctx := r.ctx
	ctx, r.cancel = context.WithCancel(r.ctx)
	go func() {
		if err := r.run(ctx, r.dist); err != nil {
			plog.Get(ctx).E().Println(err)
		}
	}()
	return
}

func (r *Runner) Run(ctx context.Context, project target.ProjectGo) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("Running project Go")
	r.ctx = ctx
	if err = r.build(ctx, project.Abs()); err != nil {
		return
	}
	r.dist = project.DistGolang()

	if err = r.Reload(); err != nil {
		return
	}

	events, err := r.watcher.Run(ctx, project.Abs())
	if err != nil {
		return
	}

	pkg := project.Manifest().Package
	for range flow.From(events).Debounce(200 * time.Millisecond) {
		if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
			log.F().Println(err)
		}
		if err = r.build(ctx, project.Abs()); err == nil {
			if err = r.Reload(); err != nil {
				log.E().Println(err)
			}
		}
		if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
			log.F().Println(err)
		}
	}
	return
}
