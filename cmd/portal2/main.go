package main

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	apphost2 "github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	apphost3 "github.com/cryptopunkscc/portal/runtime/apphost"
	portal2 "github.com/cryptopunkscc/portal/runtime/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"io"
	"os"
	"time"
)

func main() {
	d := &deps{}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	go singal.OnShutdown(d.cancel)
	if err := d.cli().Run(d.ctx); err != nil {
		return
	}
	d.cancel()
}

type deps struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (d *deps) name() string            { return "portal" }
func (d *deps) dir() string             { return target.CacheDir(d.name()) }
func (d *deps) cli() *cli.Runner        { return cli.New(d.handler()) }
func (d *deps) handler() cmd.Handler    { return portal.Handler(d.service()) }
func (d *deps) rpc() portal.Client      { return portal2.Client(d.name()) }
func (d *deps) apphost() apphost.Client { return apphost2.Full(d.ctx) }
func (d *deps) service() portal.Service {
	return &service{
		Apphost: d.apphost(),
		Rpc:     d.rpc(),
		Dir:     d.dir(),
		Run: func() error {
			return nil
		},
	}
}

type service struct {
	Apphost apphost.Client
	Rpc     portal.Client
	Run     func() error
	Dir     string
}

func (s *service) Query(ctx context.Context, query string) (err error) {
	starting := false
	if starting, err = s.startPortaldIfNeeded(); err != nil {
		return
	}
	if query == "" {
		return
	}
	if starting {
		if err = s.awaitPortaldStarted(ctx); err != nil {
			return
		}
	}
	if err = s.proceedQuery(query); err != nil {
		return
	}
	return
}

func (s *service) startPortaldIfNeeded() (starting bool, err error) {
	if fs2.IsLocked(s.Dir, "portal.lock") {
		return
	}
	if err = s.Run(); err != nil {
		return
	}
	starting = true
	return
}

func (s *service) awaitPortaldStarted(ctx context.Context) error {
	if err := flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		return apphost3.Init()
	}); err != nil {
		return err
	}
	if err := flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		return s.Rpc.Ping()
	}); err != nil {
		return err
	}
	return nil
}

func (s *service) proceedQuery(query string) error {
	conn, err := s.Apphost.Query(id.Anyone, query)
	if err != nil {
		return err
	}
	go func() { _, _ = io.Copy(conn, os.Stdin) }()
	if _, err = io.Copy(os.Stdout, conn); err != nil {
		return err
	}
	return nil
}
