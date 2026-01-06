package bundle

import (
	"context"
	"log"
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

type Repository struct {
	Apphost *apphost.Adapter
}

var _ target.Repository = &Repository{}

func (r Repository) init() {
	if r.Apphost == nil {
		r.Apphost = apphost.Default
	}
}

func (r Repository) Get(src string) (out []target.Source) {
	r.init()
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
	r.init()
	defer plog.TraceErr(&err)
	if err = r.Apphost.Connect(); err == nil {
		return
	}
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
	r.init()
	defer plog.TraceErr(&err)
	if err = r.Apphost.Connect(); err == nil {
		return
	}
	o := &Object[any]{}
	o.Resolve = target.Any[target.AppBundle[any]](Resolve_.Try)

	err = r.Apphost.Objects().Fetch(&id, o)
	out = o.AppBundle
	if out == nil && err == nil {
		err = target.ErrNotFound
	}
	return
}

func (r Repository) Scan(ctx context.Context, follow bool) (out flow.Input[Info], err error) {
	r.init()
	defer plog.TraceErr(&err)
	if err = r.Apphost.Connect(); err == nil {
		return
	}
	a := availableAppsScanner{
		log: plog.Get(ctx),
		rpc: r.Apphost.Rpc(),
		out: make(chan Info),
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

		ss, err := r.Apphost.User().Siblings(nil)
		if err != nil {
			a.log.Println(err)
			return
		}
		for id := range ss {
			a.scan(ctx, *id)
		}
	}()
	return
}

type availableAppsScanner struct {
	*apphost.Adapter
	wg  sync.WaitGroup
	log plog.Logger
	rpc rpc.Rpc
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

	ids, err := a.Objects().Scan(nil, "app.bundle", false)
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
			a.fetchInfo(*id, h)
		}
	}()
}

func (a *availableAppsScanner) fetchInfo(
	id astral.ObjectID,
	host *astral.Identity,
) {
	br := Release{}
	oc := a.Objects()
	if err := oc.Fetch(&id, &br); err != nil {
		a.log.Println(err)
		return
	}
	ma := manifest.App{}
	if err := oc.Fetch(&id, &ma); err != nil {
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
