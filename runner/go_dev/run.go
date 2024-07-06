package go_dev

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/target"
	"time"
)

type Runner struct {
	watcher *golang.Watcher
	run     target.Run[target.DistExec]
	send    target.MsgSend
	dist    target.DistExec
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewRunner(
	run target.Run[target.DistExec],
	send target.MsgSend,
) target.Runner[target.ProjectGo] {
	return &Runner{
		watcher: golang.NewWatcher(),
		send:    send,
		run:     run,
	}
}

func NewAdapter(run target.Run[target.Portal]) func(
	_ target.NewApi,
	send target.MsgSend,
) target.Runner[target.ProjectGo] {
	return func(newApi target.NewApi, send target.MsgSend) target.Runner[target.ProjectGo] {
		run := func(ctx context.Context, src target.DistExec) (err error) {
			newApi(ctx, src) // initiate connection
			return run(ctx, src)
		}
		return NewRunner(run, send)
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
	build := dist.NewGoRunner().Run
	if err = build(ctx, project); err != nil {
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
			log.E().Println(err)
		}
		if err = build(ctx, project); err == nil {
			if err = r.Reload(); err != nil {
				log.E().Println(err)
			}
		}
		time.Sleep(2 * time.Second)
		if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
			log.E().Println(err)
		}
	}
	return
}
