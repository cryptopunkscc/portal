package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/token"
)

var DistRun target.Run[target.DistExec] = dist

func dist(ctx context.Context, src target.DistExec, args ...string) (err error) {
	abs := src.Target().Executable().Abs()
	t, err := token.Repository{}.Get(src.Manifest().Package)
	if err != nil {
		return err
	}
	return Cmd{}.Run(ctx, t.Token.String(), abs, args...)
}
