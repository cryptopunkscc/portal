package tray

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"log"
)
import "github.com/getlantern/systray"

func Run(ctx context.Context) {
	t := tray{}
	systray.SetTitle(portal.Name)
	launcherItem := systray.AddMenuItem("Launcher", "Launcher")
	go onMenuItemClick(launcherItem, func() {
		go func() {
			if err := portal.Open(ctx, "launcher").Run(); err != nil {
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

type tray struct{}

func (t *tray) onReady() {
	log.Println("portal tray start")
}

func (t *tray) onExit() {
	log.Println("portal tray exit")
}

func onMenuItemClick(item *systray.MenuItem, onClick func()) {
	for range item.ClickedCh {
		onClick()
	}
}
