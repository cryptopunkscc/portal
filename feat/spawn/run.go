package spawn

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func Run[T target.Portal](
	ctx context.Context,
	resolve runtime.Resolve[T],
	action string,
	src string,
) (err error) {

	// resolve apps from given source
	apps, err := resolve(src)
	if len(apps) == 0 {
		return errors.Join(fmt.Errorf("no apps found in %s", src), err)
	}

	// execute multiple targets as separate processes
	if len(apps) > 1 {
		return portal.Spawn(nil, ctx, apps, action)
	}
	return
}
