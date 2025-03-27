package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/fsnotify/fsnotify"
)

func (s *Service[T]) ObserveApps(ctx context.Context, opts ListAppsOpts) (out <-chan target.Portal_, err error) {
	log := plog.Get(ctx)
	log.Println("Observing...")

	watch, err := fs2.NotifyWatch(ctx, s.Config.Apps, 0)
	if err != nil {
		return
	}

	results := make(chan target.Portal_)
	out = results
	go func() {
		defer close(results)
		resolve := apps.Resolver[target.App_]()

		for _, bundle := range resolve.List(s.apps()) {
			if opts.Hidden || !bundle.Manifest().Hidden {
				results <- bundle
			}
		}

		for event := range watch {
			log.Println("Event:", event)
			if event.Op != fsnotify.Write {
				continue
			}
			if file, err := source.File(event.Name); err == nil {
				for _, bundle := range resolve.List(file) {
					if opts.Hidden || !bundle.Manifest().Hidden {
						results <- bundle
					}
					break
				}
			}
		}
	}()
	return
}
