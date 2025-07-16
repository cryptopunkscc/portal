package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (s *Service) InstallApp(ctx context.Context, src string) (out <-chan any, err error) {
	plog.TraceErr(&err)
	list := s.Installer().Dispatcher().Provide(src)
	if list == nil {
		err = target.ErrNotFound
		return
	}
	c := make(chan any)
	out = c
	go func() {
		defer close(c)
		for _, r := range list {
			if ctx.Err() != nil {
				break
			}
			if err = r.Run(ctx); err != nil {
				c <- err
			} else {
				c <- r.Manifest()
			}
		}
	}()
	return
}
