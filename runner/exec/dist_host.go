package exec

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/source"
)

func (r Runner) DistHost() *target.SourceRunner[target.Portal_] {
	return &target.SourceRunner[target.Portal_]{
		Resolve: target.Any[target.Portal_](target.Try(dist.Resolve_)),
		Runner:  &DistHostRunner{DistRunner{r}},
	}
}

type DistHostRunner struct{ DistRunner }

func (r *DistHostRunner) Run(ctx context.Context, src target.Portal_, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if src.Manifest().Runtime == "" {
		return errors.New("DistHostRunner requires a schema declared in manifest")
	}

	log := plog.Get(ctx).Type(r)
	repo := target.SourcesRepository[target.DistExec]{
		Sources: []target.Source{source.Dir(r.Apps)},
		Resolve: target.Any[target.DistExec](exec.ResolveDist.Try),
	}
	hostId := src.Manifest().Runtime
	opt := portald.OpenOpt{}
	if opt.Load(ctx); len(opt.Schema) > 0 {
		hostId = opt.Schema + "." + hostId
	}

	log.Println("running:", hostId, src.Manifest().Package, args)

	hostDist := repo.First(hostId)
	if hostDist == nil {
		return target.ErrNotFound
	}

	execFile := hostDist.Runtime().Target().Exec
	execFile = filepath.Join(hostDist.Abs(), execFile)
	execFile = strconv.Quote(execFile)

	args = slices.Insert(args, 0, strconv.Quote(src.Abs()))
	return r.RunApp(ctx, *src.Manifest(), execFile, args...)
}
