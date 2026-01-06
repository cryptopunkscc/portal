package app

import (
	"context"
	"io/fs"
	"log"
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/source"
)

type Objects struct {
	*apphost.Adapter
}

var _ source.Provider = &Objects{}

func (r Objects) GetSource(src string) (out source.Source) {
	out, _ = r.GetAppBundle(src)
	return
}

func (r Objects) GetAppBundle(src string) (out *Bundle, err error) {
	id, err := astral.ParseID(src)
	if err == nil {
		return r.GetByObjectID(*id)
	}
	return r.GetByNameOrPkg(src)
}

func (r Objects) GetByNameOrPkg(name string) (out *Bundle, err error) {
	defer plog.TraceErr(&err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for info := range r.Scan(ctx, false) {
		if !info.Manifest.Match(name) {
			continue
		}
		if info.Host == nil {
			info.Host = r.HostID()
		}
		out, err = r.GetByObjectID(*info.BundleID, *info.Host)
		return
	}
	err = fs.ErrNotExist
	return
}

func (r Objects) GetByObjectID(id astral.ObjectID, host ...astral.Identity) (out *Bundle, err error) {
	defer plog.TraceErr(&err)
	var obj astral.Object
	for _, identity := range host {
		r.Client = r.WithTarget(&identity)
		if obj, err = r.Objects().Get(&id); err == nil {
			if out, _ = obj.(*Bundle); out != nil {
				return
			}
		}
		return
	}
	return
}

func (r Objects) Scan(ctx context.Context, follow bool) (out flow.Input[ReleaseInfo]) {
	a := availableAppsScanner{
		Adapter: r.Adapter,
		log:     plog.Get(ctx),
		out:     make(chan ReleaseInfo),
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

		if c, err := r.User().Siblings(astral.NewContext(ctx)); err == nil {
			for id := range c {
				a.scan(ctx, *id)
			}
		}
	}()
	return
}

type availableAppsScanner struct {
	*apphost.Adapter
	wg  sync.WaitGroup
	log plog.Logger
	rpc rpc.Rpc
	out chan ReleaseInfo
}

type ReleaseInfo struct {
	Manifest Manifest
	ReleaseMetadata
	ReleaseID *astral.ObjectID
	Host      *astral.Identity
}

func (a *availableAppsScanner) scan(ctx context.Context, host ...astral.Identity) {
	if ctx.Err() != nil {
		return
	}

	var h *astral.Identity
	if len(host) > 0 {
		h = &host[0]
	}

	ids, err := a.Objects().Scan(astral.NewContext(ctx), "", false)
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
	metadata := &ReleaseMetadata{}
	oc := a.Objects()

	o, err := oc.Get(&id)
	if err != nil {
		return
	}
	metadata, ok := o.(*ReleaseMetadata)
	if !ok {
		return
	}

	//if err := oc.Fetch(&id, &metadata); err != nil {
	//	a.log.Println(err)
	//	return
	//}

	manifest := Manifest{}
	if err := oc.Fetch(metadata.ManifestID, &manifest); err != nil {
		a.log.Println(err)
		return
	}
	i := ReleaseInfo{
		Manifest:        manifest,
		ReleaseMetadata: *metadata,
		ReleaseID:       &id,
		Host:            host,
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
