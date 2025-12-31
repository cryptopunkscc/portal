package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Opt struct {
	Token string `cli:"token t"`
	Query string `cli:"query q"`
	Open  bool   `cli:"open o"`
	Dev   bool   `cli:"dev d"`
	Order string `cli:"order"`
}

func (a *Application) Run(ctx context.Context, opt Opt, cmd ...string) (err error) {
	defer plog.TraceErr(&err)
	log := plog.Get(ctx).Type(a).Set(&ctx)
	if os.Getenv("ENABLE_PORTAL_APPHOST_LOG") == "true" {
		a.Apphost.Log = log
	}
	if err = a.Configure(); err != nil {
		if err = a.handleConfigurationError(ctx, err); err != nil {
			return
		}
	}
	if len(opt.Token) > 0 {
		a.Apphost.Token = opt.Token
		if err = a.Apphost.Reconnect(); err != nil {
			return
		}
	}
	opt.Open = opt.Open || opt.Query != ""
	if err = a.portald().Ping(); err != nil {
		if err = a.startPortald(ctx); err != nil {
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
			o := apphost.OpenOpt{
				App:  cmd[0],
				Args: strings.Join(cmd[1:], " "),
			}
			if opt.Open {
				err = a.startApp(ctx, o)
			} else {
				err = a.runApp(ctx, o)
			}
		}()
	}
	wg.Wait()
	log.Println("exit")
	return
}
