package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cryptopunkscc/portal/api/portald"
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
		a.Apphost.AuthToken = opt.Token
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
			o := &portald.OpenOpt{}
			if opt.Dev {
				o.Schema = "dev"
				o.Order = []int{2, 1, 0}
			}
			if opt.Order != "" {
				o.Order = nil
				i := 0
				for _, s := range strings.Split(opt.Order, ",") {
					if i, err = strconv.Atoi(s); err != nil {
						return
					}
					o.Order = append(o.Order, i)
				}
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
