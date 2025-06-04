package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target/bundle"
	"log"
	"sync"
)

type AvailableApp struct {
	manifest.App
	bundle.Release
	ReleaseID *astral.ObjectID
}

func (s *Service[T]) AvailableApps(ctx context.Context, follow bool) (out flow.Input[AvailableApp], err error) {
	a := availableAppsScanner{}
	a.log = plog.Get(ctx)
	a.rpc = s.Apphost.Rpc()
	a.arg = objects.ScanArgs{
		Type:   bundle.Release{}.ObjectType(),
		Zone:   astral.ZoneAll,
		Follow: follow,
	}
	a.out = make(chan AvailableApp)
	a.wg = &sync.WaitGroup{}
	out = a.out

	a.wg.Add(1)
	go func() {
		a.wg.Wait()
		close(a.out)
	}()
	go func() {
		defer a.wg.Done()

		// local apps
		a.scan()

		// remote apps
		uc := user.Client{Rpc: a.rpc}
		ss, err := uc.Siblings()
		if err != nil {
			log.Println(err)
			return
		}
		for id := range ss {
			a.scan(id.String())
		}
	}()
	return
}

type availableAppsScanner struct {
	log plog.Logger
	rpc rpc.Rpc
	wg  *sync.WaitGroup
	arg objects.ScanArgs
	out chan AvailableApp
}

func (a availableAppsScanner) scan(target ...string) {
	oc := objects.Client(a.rpc, target...)
	ids, err := oc.Scan(a.arg)
	if err != nil {
		log.Println(err)
		return
	}
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for id := range ids {

			br := bundle.Release{}
			oa := objects.ReadArgs{
				ID:   id,
				Zone: astral.ZoneAll,
			}
			if err := oc.Fetch(oa, &br); err != nil {
				log.Println(err)
				continue
			}

			oa.ID = *br.ManifestID
			ma := manifest.App{}
			if err := oc.Fetch(oa, &ma); err != nil {
				log.Println(err)
				continue
			}

			aa := AvailableApp{
				App:       ma,
				Release:   br,
				ReleaseID: &id,
			}
			a.out <- aa
		}
	}()
}
