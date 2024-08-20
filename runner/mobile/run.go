package mobile

import (
	"context"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
)

func Run[T target.Portal_](request mobile.Request) target.Run[T] {
	return func(ctx context.Context, src T) (err error) {
		return request(src.Abs())
	}
}
