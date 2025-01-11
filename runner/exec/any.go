package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"slices"
)

func RunAny(runner func(string) string) target.Run[target.Portal_] {
	return func(ctx context.Context, portal target.Portal_, args ...string) (err error) {
		schema := portal.Manifest().Schema
		src := runner(schema)
		if src == "" {
			return fmt.Errorf("unknown schema %v", schema)
		}
		args = slices.Insert(args, 0, portal.Abs())
		return RunCmd(ctx, src, args...)
	}
}

func AnyRun(cacheDir string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		log := plog.Get(ctx)
		schema := src.Manifest().Schema
		log.Println("run schema:", schema)
		runners, err := target.FindByPath(source.File, exec2.ResolveBundle).ById(appstore.Path).Call(ctx, schema)
		if err != nil {
			return
		}
		log.Println("run runners:", runners)
		if len(runners) == 0 {
			return target.ErrNotFound
		}
		runner := runners[0]
		log.Println("running:", runner.Manifest().Name)
		args = slices.Insert(args, 0, src.Abs())
		err = BundleRun(cacheDir).Call(ctx, runner, args...)
		log.Println("running err:", err)
		return
	}
}
