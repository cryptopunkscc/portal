package go_project

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/flow"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/reload"
	golang2 "github.com/cryptopunkscc/portal/target/go"
	"time"
)

func Runner() *target.SourceRunner[target.ProjectGo] {
	return &target.SourceRunner[target.ProjectGo]{
		Resolve: target.Any[target.ProjectGo](golang2.ResolveProject.Try),
		Runner: &ReRunner{
			watcher: golang.NewWatcher(),
			run:     exec.DefaultRunner().Dist().Runner.Run,
		},
	}
}

type ReRunner struct {
	watcher *golang.Watcher
	run     target.Run[target.DistExec]
	send    target.MsgSend
	dist    target.DistExec
	ctx     context.Context
	cancel  context.CancelFunc
	args    []string
}

func (r *ReRunner) Reload() (err error) {
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

func (r *ReRunner) Run(ctx context.Context, project target.ProjectGo, args ...string) (err error) {
	r.args = args
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("Running project Go")
	if err = deps.RequireBinary("go"); err != nil {
		return
	}
	r.ctx = ctx
	build := golang2.BuildProject()
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
	r.send = reload.Start(ctx, project, r.Reload, nil)
	for range flow.From(events).Debounce(200 * time.Millisecond) {
		if err := r.sendMsg(pkg, target.DevChanged); err != nil {
			log.E().Println(err)
		}
		if err = build(ctx, project); err == nil {
			if err = r.Reload(); err != nil {
				log.E().Println(err)
			}
		}
		time.Sleep(2 * time.Second)
		if err := r.sendMsg(pkg, target.DevRefreshed); err != nil {
			log.E().Println(err)
		}
	}
	return
}

func (r *ReRunner) sendMsg(pkg string, event target.Event) (err error) {
	if r.send != nil {
		return r.send(target.NewMsg(pkg, event))
	}
	return
}
