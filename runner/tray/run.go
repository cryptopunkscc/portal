package tray

import (
	"context"
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
)
import "github.com/getlantern/systray"

type Runner struct {
	open target.Dispatch
}

func NewRunner(open target.Dispatch) target.Tray {
	return (&Runner{open: open}).Run
}

func (t *Runner) Run(ctx context.Context) {
	systray.SetTitle(portal.Name)
	launcherItem := systray.AddMenuItem("Launcher", "Launcher")
	go onMenuItemClick(launcherItem, func() {
		go func() {
			if err := t.open(ctx, "launcher"); err != nil {
				log.Println("launcher:", err)
			}
		}()
	})
	quit := systray.AddMenuItem("Quit ", "Quit")
	go onMenuItemClick(quit, func() {
		systray.Quit()
		if err := exec.Shutdown(); err != nil {
			log.Println("quit:", err)
		}
	})
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()
	systray.Run(t.onReady, t.onExit)
}

func (t *Runner) onReady() {
	log.Println("portal tray start")
}

func (t *Runner) onExit() {
	log.Println("portal tray exit")
}

func onMenuItemClick(item *systray.MenuItem, onClick func()) {
	for range item.ClickedCh {
		onClick()
	}
}
