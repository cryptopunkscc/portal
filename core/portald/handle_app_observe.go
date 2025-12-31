package portald

import (
	"context"

	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/fsnotify/fsnotify"
)

func (s *Service) ObserveApps(ctx context.Context, opts ListAppsOpts) (out <-chan ObservedApp, err error) {
	log := plog.Get(ctx)
	log.Println("Observing...")

	watch, err := fs2.NotifyWatch(ctx, s.Config.Apps, 0)
	if err != nil {
		return
	}

	installed := sig.Map[string, bool]{}
	results := make(chan ObservedApp)
	out = results

	go func() {
		<-ctx.Done()
		close(results)
	}()

	// list installed apps
	go func() {
		for _, a := range source.CollectT[app.App](s.appsRef(), &app.Dist{}, &app.Bundle{}) {
			metadata := a.GetDist().Metadata
			if opts.includes(metadata) {
				installed.Set(metadata.Package, true)
				results <- ObservedApp{
					Manifest:  metadata.Manifest,
					Installed: true,
				}
			}
		}

		// observe installed apps
		go func() {
			for event := range watch {
				if event.Op != fsnotify.Create && event.Op != fsnotify.Write {
					continue
				}
				for _, a := range source.CollectT[app.App](s.appsRef(), &app.Dist{}, &app.Bundle{}) {
					metadata := a.GetDist().Metadata
					if opts.includes(metadata) {
						installed.Set(metadata.Package, true)
						results <- ObservedApp{
							Manifest:  metadata.Manifest,
							Installed: true,
						}
					}
				}
			}
		}()

		// observe available apps
		go func() {
			for info := range s.AvailableApps(ctx, true) {
				log.Println("new available app:", info.Manifest)
				_, b := installed.Get(info.Manifest.Package)
				results <- ObservedApp{
					Manifest:  info.Manifest,
					Installed: b,
				}
			}
		}()
	}()
	return
}

type ObservedApp struct {
	app.Manifest
	Installed bool `json:"installed"`
}
