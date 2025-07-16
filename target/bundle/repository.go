package bundle

import (
	"context"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"log"
	"sync"
)

type Repository struct {
	Apphost apphost.Client
}

var _ target.Repository = &Repository{}

func (r Repository) Get(src string) (out []target.Source) {
	id, err := astral.ParseID(src)
	var b target.Bundle_
	if err == nil {
		b, err = r.GetByObjectID(*id)
	} else {
		b, err = r.GetByNameOrPkg(src)
	}
	if err == nil {
		out = append(out, b)
	}
	return
}

func (r Repository) GetByNameOrPkg(name string) (out target.Bundle_, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	infos, err := r.Scan(ctx, false)
	if err != nil {
		return
	}
	for info := range infos {
		if !info.Manifest.Match(name) {
			continue
		}
		var host []astral.Identity
		if info.Host != nil {
			host = append(host, *info.Host)
		}
		out, err = r.GetByObjectID(*info.Release.BundleID, host...)
		return
	}
	err = target.ErrNotFound
	return
}

func (r Repository) GetByObjectID(id astral.ObjectID, host ...astral.Identity) (out target.Bundle_, err error) {
	defer plog.TraceErr(&err)
	c := objects.Op(r.Apphost.Rpc(), identities(host).Strings()...)
	o := &Object[any]{}
	o.Resolve = target.Any[target.AppBundle[any]](Resolve_.Try)
	err = c.Fetch(objects.ReadArgs{ID: id}, o)
	out = o.AppBundle
	if out == nil && err == nil {
		err = target.ErrNotFound
	}
	return
}

func (r Repository) Scan(ctx context.Context, follow bool) (out flow.Input[Info], err error) {
	a := availableAppsScanner{
		log: plog.Get(ctx),
		rpc: r.Apphost.Rpc(),
		out: make(chan Info),
		arg: objects.ScanArgs{
			Type:   Release{}.ObjectType(),
			Zone:   astral.ZoneAll,
			Follow: follow,
		},
	}

	out = a.out

	a.wg.Add(1)
	go func() {
		a.wg.Wait()
		close(a.out)
	}()
	go func() {
		defer a.wg.Done()

		// local apps
		a.scan(ctx)

		if ctx.Err() != nil {
			return
		}

		// remote apps
		uc := user.Op(r.Apphost)
		ss, err := uc.Siblings()
		if err != nil {
			a.log.Println(err)
			return
		}
		for id := range ss {
			a.scan(ctx, id)
		}
	}()
	return
}

type availableAppsScanner struct {
	wg  sync.WaitGroup
	log plog.Logger
	rpc rpc.Rpc
	arg objects.ScanArgs
	out chan Info
}

func (a *availableAppsScanner) scan(ctx context.Context, host ...astral.Identity) {
	if ctx.Err() != nil {
		return
	}

	var h *astral.Identity
	if len(host) > 0 {
		h = &host[0]
	}

	oc := objects.Op(a.rpc, identities(host).Strings()...)
	ids, err := oc.Scan(a.arg)
	if err != nil {
		log.Println(err)
		return
	}
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for id := range ids {
			if ctx.Err() != nil {
				continue
			}
			a.fetchInfo(oc, id, h)
		}
	}()
}

func (a *availableAppsScanner) fetchInfo(
	oc objects.OpClient,
	id astral.ObjectID,
	host *astral.Identity,
) {
	br := Release{}
	oa := objects.ReadArgs{
		ID:   id,
		Zone: astral.ZoneAll,
	}
	if err := oc.Fetch(oa, &br); err != nil {
		a.log.Println(err)
		return
	}

	oa.ID = *br.ManifestID
	ma := manifest.App{}
	if err := oc.Fetch(oa, &ma); err != nil {
		a.log.Println(err)
		return
	}
	i := Info{
		Manifest:  ma,
		Release:   br,
		ReleaseID: &id,
		Host:      host,
	}
	a.out <- i
}

type identities []astral.Identity

func (i identities) Strings() (out []string) {
	for _, i := range i {
		out = append(out, i.String())
	}
	return
}
