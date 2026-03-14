package main

import (
	"context"
	"errors"

	"github.com/cryptopunkscc/portal/pkg/client"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"github.com/getlantern/systray"
)

type Tray struct {
	Portald client.Portald
	log     plog.Logger
}

func (t *Tray) Run(ctx context.Context) (err error) {
	if t.Portald.Astrald == nil {
		t.Portald = client.Default.Portald()
	}

	if err = t.Portald.Ping(); err != nil {
		return errors.New("portal-tray requires portal-app running")
	}

	t.log = plog.Get(ctx).Type(t).Set(&ctx)
	t.Portald.Log = t.log

	go func() {
		t.Portald.Join()
		systray.Quit()
	}()
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()

	systray.Run(t.onReady, t.onExit)
	return
}

func (t *Tray) onReady() {
	t.log.Println("portal tray start")
	launcher := systray.AddMenuItem("Launcher", "Launcher")
	quit := systray.AddMenuItem("Quit ", "Quit")

	go func() {
		for {
			select {
			case <-launcher.ClickedCh:
				go func() {
					if err := t.Portald.Open(client.OpenOpt{App: "launcher"}); err != nil {
						t.log.Println("launcher:", err)
					}
				}()
			case <-quit.ClickedCh:
				if err := t.Portald.Close(); err != nil {
					t.log.Println("quit:", err)
					systray.Quit()
				}
			}
		}
	}()
}

func (t *Tray) onExit() {
	t.log.Println("portal tray exit")
}
