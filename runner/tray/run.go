package tray

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"github.com/getlantern/systray"
)

type Api interface {
	Open(src string) error
	Close() error
	Ping() error
	Await()
}

type AwaitClose func()

func NewRun(api Api) target.Tray {
	return (&Runner{api: api}).Run
}

type Runner struct {
	api Api
	log plog.Logger
}

func (t *Runner) Run(ctx context.Context) (err error) {
	if err = t.api.Ping(); err != nil {
		return errors.New("portal-tray requires portal-app running")
	}

	t.log = plog.Get(ctx).Type(t).Set(&ctx)

	go func() {
		t.api.Await()
		systray.Quit()
	}()
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()

	systray.Run(t.onReady, t.onExit)
	return
}

func (t *Runner) onReady() {
	t.log.Println("portal tray start")
	systray.SetTitle(portal.Name)
	launcher := systray.AddMenuItem("Launcher", "Launcher")
	quit := systray.AddMenuItem("Quit ", "Quit")

	go func() {
		for {
			select {
			case <-launcher.ClickedCh:
				go func() {
					if err := t.api.Open("launcher"); err != nil {
						t.log.Println("launcher:", err)
					}
				}()
			case <-quit.ClickedCh:
				if err := t.api.Close(); err != nil {
					t.log.Println("quit:", err)
					systray.Quit()
				}
			}
		}
	}()
}

func (t *Runner) onExit() {
	t.log.Println("portal tray exit")
}
