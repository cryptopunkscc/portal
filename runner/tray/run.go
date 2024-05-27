package tray

import (
	"context"
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/getlantern/systray"
)

func NewRun(open target.Dispatch) target.Tray {
	return (&Runner{open: open}).Run
}

type Runner struct {
	open target.Dispatch
	log  plog.Logger
}

func (t *Runner) Run(ctx context.Context) {
	t.log = plog.Get(ctx).Type(t).Set(&ctx)
	systray.SetTitle(portal.Name)
	launcherItem := systray.AddMenuItem("Launcher", "Launcher")
	go onMenuItemClick(launcherItem, func() {
		go func() {
			if err := t.open(ctx, "launcher"); err != nil {
				t.log.Println("launcher:", err)
			}
		}()
	})
	quit := systray.AddMenuItem("Quit ", "Quit")
	go onMenuItemClick(quit, func() {
		systray.Quit()
		if err := exec.Shutdown(); err != nil {
			t.log.Println("quit:", err)
		}
	})
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()
	systray.Run(t.onReady, t.onExit)
}

func (t *Runner) onReady() {
	t.log.Println("portal tray start")
}

func (t *Runner) onExit() {
	t.log.Println("portal tray exit")
}

func onMenuItemClick(item *systray.MenuItem, onClick func()) {
	for range item.ClickedCh {
		onClick()
	}
}
