package portald

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

type SetupOpts struct {
	User string `query:"user u" cli:"user u"`
	Apps string `query:"apps a" cli:"apps a"`
}

func (s *Service) Setup(ctx context.Context, opts SetupOpts) (err error) {
	defer plog.PrintTrace(&err)
	if s.HasUser() {
		return plog.Errorf("setup already performed")
	}
	if len(opts.User) > 0 {
		if err = s.CreateUser(opts.User); err != nil {
			return
		}
	}
	if len(opts.Apps) > 0 {
		for _, install := range s.Installer().Dispatcher().Provide(opts.Apps) {
			if err = install.Run(ctx); err != nil {
				return
			}
		}
	}
	return
}
