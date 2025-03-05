package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runtime/tokens"
)

func DistRunner() target.Run[target.DistExec] {
	return func(ctx context.Context, src target.DistExec, args ...string) (err error) {
		abs := src.Target().Executable().Abs()
		token, err := tokens.Repository{}.Get(src.Manifest().Package)
		if err != nil {
			return err
		}
		return RunCmd(ctx, token.Token.String(), abs, args...)
	}
}
