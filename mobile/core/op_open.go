package core

import (
	"io/fs"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/ops"
	bind "github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/runner/goja"
	"github.com/cryptopunkscc/portal/pkg/source"
)

func (srv *Service) OpOpen(ctx *astral.Context, q *ops.Query, args OpOpenArgs) (err error) {
	ch := q.AcceptChannel(channel.WithFormats(args.In, args.Out))
	defer func() {
		if err != nil {
			err = ch.Send(astral.Err(err))
		} else {
			err = ch.Send(&astral.Ack{})
		}
		_ = ch.Close()
	}()

	app, err := srv.app(args.App)
	if err != nil {
		return
	}
	for _, s := range source.Collect(app,
		&goja.BundleRunner{},
		&htmlBundleRunner{},
	) {
		switch r := s.(type) {
		case *goja.BundleRunner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Start(ctx)
		case *htmlBundleRunner:
			return r.Start()
		}
	}

	return fs.ErrInvalid
}

type OpOpenArgs struct {
	App string
	In  string `query:"optional"`
	Out string `query:"optional"`
}
