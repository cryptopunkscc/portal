package main

import (
	"context"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"sync"
	"time"
)

type Opt struct {
	Query string `cli:"query q"`
	Open  bool   `cli:"open o"`
	Dev   bool   `cli:"dev d"`
}

func (a Application) Run(ctx context.Context, opt Opt, cmd ...string) (err error) {
	opt.Open = opt.Open || opt.Query != ""
	log := plog.Get(ctx).Type(a).Set(&ctx)
	if err = a.Connect(ctx); err != nil {
		return
	}
	a.Portal.Logger(log)
	if err = a.Portal.Ping(); err != nil {
		if err = startPortald(ctx, a.Portal); err != nil {
			return
		}
	}
	wg := sync.WaitGroup{}
	if opt.Query != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(cmd) > 0 {
				time.Sleep(200 * time.Millisecond)
			}
			if ctx.Err() != nil {
				return
			}
			err = a.queryApp(ctx, opt.Query)
		}()
	}
	if len(cmd) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd = fixCmd(cmd)
			o := &portald.OpenOpt{}
			if opt.Dev {
				o.Schema = "dev"
				o.Order = []int{2, 1, 0}
			}
			if opt.Open {
				err = a.startApp(ctx, o, cmd)
			} else {
				err = a.runApp(ctx, o, cmd)
			}
		}()
	}
	wg.Wait()
	log.Println("exit")
	return
}
