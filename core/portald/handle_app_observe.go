package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/app"
	"github.com/cryptopunkscc/portal/target/source"
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
	resolve := app.Resolve_

	go func() {
		<-ctx.Done()
		close(results)
	}()

	// list installed apps
	go func() {
		for _, bundle := range resolve.List(s.apps()) {
			if opts.Hidden || !bundle.Config().Hidden {
				installed.Set(bundle.Manifest().Package, true)
				results <- ObservedApp{
					App:       *bundle.Manifest(),
					Installed: true,
				}
			}
		}

		// observe installed apps
		go func() {
			for event := range watch {
				log.Println("Event:", event)
				if event.Op != fsnotify.Create && event.Op != fsnotify.Write {
					continue
				}
				if file, err := source.File(event.Name); err == nil {
					log.Println("new installed file:", file.Abs())
					for _, bundle := range resolve.List(file) {
						log.Println("new installed app:", *bundle.Manifest())
						if opts.Hidden || !bundle.Config().Hidden {
							installed.Set(bundle.Manifest().Package, true)
							log.Println("new installed app sending:", *bundle.Manifest())
							results <- ObservedApp{
								App:       *bundle.Manifest(),
								Installed: true,
							}
						}
						break
					}
				}
			}
		}()

		// observe available apps
		go func() {
			if apps, err := s.AvailableApps(ctx, true); err == nil {
				for info := range apps {
					log.Println("new available app:", info.Manifest)
					_, b := installed.Get(info.Manifest.Package)
					results <- ObservedApp{
						App:       info.Manifest,
						Installed: b,
					}
				}
			} else {
				log.Println("Error observing apps:", err)
			}
		}()
	}()
	return
}

type ObservedApp struct {
	manifest.App
	Installed bool `json:"installed"`
}
