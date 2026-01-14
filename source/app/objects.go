package app

import (
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

func (r Objects) Default() Objects {
	r.Adapter = apphost.Default
	return r
}

func (r Objects) GetSource(src string) (out source.Source) {
	if b, err := r.GetAppBundle(astral.NewContext(nil), src); err == nil {
		return b
	}
	return
}

func (r Objects) GetAppBundle(ctx *astral.Context, src string) (out *Bundle, err error) {
	id, err := astral.ParseID(src)
	if err == nil {
		return r.GetByObjectID(ctx, *id)
	}
	return r.GetByNameOrPkg(ctx, src)
}

func (r Objects) GetByNameOrPkg(ctx *astral.Context, name string) (out *Bundle, err error) {
	defer plog.TraceErr(&err)
	ctx, cancel := ctx.WithCancel()
	defer cancel()
	for info := range r.Scan(ctx, false) {
		if !info.Manifest.Match(name) {
			continue
		}
		if info.Host == nil {
			info.Host = r.HostID()
		}
		out, err = r.GetByObjectID(ctx, *info.BundleID, *info.Host)
		return
	}
	return nil, fs.ErrNotExist
}

func (r Objects) GetByObjectID(ctx *astral.Context, id astral.ObjectID, host ...astral.Identity) (out *Bundle, err error) {
	defer plog.TraceErr(&err)
	var obj astral.Object
	if len(host) == 0 {
		host = append(host, *r.HostID())
	}
	for _, identity := range host {
		r.Client = r.WithTarget(&identity)
		if obj, err = r.Objects().Get(ctx, &id); err == nil {
			switch v := obj.(type) {
			case *Bundle:
				out = v
				return
			case *ReleaseMetadata:
				if obj, err = r.Objects().Get(ctx, v.BundleID); err != nil {
					continue
				}
				if out, _ = obj.(*Bundle); out != nil {
					return
				}
			}
		}
	}
	return nil, fs.ErrNotExist
}

func (r Objects) Scan(ctx *astral.Context, follow bool) (out flow.Input[ReleaseInfo]) {
	a := availableAppsScanner{
		Adapter: r.Adapter,
		log:     plog.Get(ctx),
		follow:  follow,
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

		if c, err := r.User().Siblings(ctx); err == nil {
			for id := range c {
				a.scan(ctx, *id)
			}
		}
	}()
	return
}

type availableAppsScanner struct {
	*apphost.Adapter
	wg     sync.WaitGroup
	log    plog.Logger
	rpc    rpc.Rpc
	follow bool
	out    chan ReleaseInfo
}

type ReleaseInfo struct {
	Manifest Manifest
	ReleaseMetadata
	ReleaseID *astral.ObjectID
	Host      *astral.Identity
}

func (a *availableAppsScanner) scan(ctx *astral.Context, host ...astral.Identity) {
	if ctx.Err() != nil {
		return
	}

	var h *astral.Identity
	if len(host) > 0 {
		h = &host[0]
	}

	ids, err := a.Objects().Scan(ctx, "main", a.follow)
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
			a.fetchInfo(ctx, *id, h)
		}
	}()
}

func (a *availableAppsScanner) fetchInfo(
	ctx *astral.Context,
	id astral.ObjectID,
	host *astral.Identity,
) {
	metadata := &ReleaseMetadata{}
	oc := a.Objects()

	o, err := oc.Get(ctx, &id)
	if err != nil {
		return
	}
	metadata, ok := o.(*ReleaseMetadata)
	if !ok {
		return
	}

	if o, err = oc.Get(ctx, metadata.ManifestID); err != nil {
		a.log.Println(err)
		return
	}
	manifest, ok := o.(*Manifest)
	if !ok {
		return
	}
	i := ReleaseInfo{
		Manifest:        *manifest,
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
