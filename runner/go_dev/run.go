package go_dev

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/flow"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"time"
)

type runner struct {
	watcher *golang.Watcher
	run     target.Run[target.DistExec]
	send    target.MsgSend
	dist    target.DistExec
	ctx     context.Context
	cancel  context.CancelFunc
	args    []string
}

func Adapter(run target.Run[target.DistExec]) func(
	_ bind.NewCore,
	send target.MsgSend,
) target.ReRunner[target.ProjectGo] {
	return func(newCore bind.NewCore, send target.MsgSend) target.ReRunner[target.ProjectGo] {
		run := func(ctx context.Context, src target.DistExec, args ...string) (err error) {
			newCore(ctx, src) // initiate connection
			return run(ctx, src, args...)
		}
		return ReRunner(run, send)
	}
}

func ReRunner(
	run target.Run[target.DistExec],
	send target.MsgSend,
) target.ReRunner[target.ProjectGo] {
	return &runner{
		watcher: golang.NewWatcher(),
		send:    send,
		run:     run,
	}
}

func (r *runner) Reload() (err error) {
	if r.cancel != nil {
		r.cancel()
	}
	ctx := r.ctx
	ctx, r.cancel = context.WithCancel(r.ctx)
	go func() {
		if err := r.run(ctx, r.dist, r.args...); err != nil {
			plog.Get(ctx).E().Println("reload", err)
		}
	}()
	return
}

func (r *runner) Run(ctx context.Context, project target.ProjectGo, args ...string) (err error) {
	r.args = args
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("Running project Go")
	if err = deps.RequireBinary("go"); err != nil {
		return
	}
	r.ctx = ctx
	build := go_build.Runner()
	if err = build(ctx, project); err != nil {
		return
	}
	r.dist = project.Dist()

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
