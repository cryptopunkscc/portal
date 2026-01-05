package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"reflect"
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/source"
)

type Objects struct {
	astrald.Client
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
		var host []astral.Identity
		if info.Host != nil {
			host = append(host, *info.Host)
		}
		out, err = r.GetByObjectID(*info.BundleID, host...)
		return
	}
	err = fs.ErrNotExist
	return
}

func (r Objects) GetByObjectID(id astral.ObjectID, host ...astral.Identity) (out *Bundle, err error) {
	defer plog.TraceErr(&err)
	var c *channel.Channel
	for _, identity := range host {
		out, _, err = func() (o *Bundle, n int64, err error) {
			c, err = r.QueryChannel(identity.String(), "objects.read", query.Args{"id": id.String()})
			if err != nil {
				return
			}
			defer c.Close()
			return astralRead[*Bundle](c.Transport())
		}()
		return
	}
	return
}

func (r Objects) Scan(ctx context.Context, follow bool) (out flow.Input[ReleaseInfo]) {
	a := availableAppsScanner{
		log: plog.Get(ctx),
		out: make(chan ReleaseInfo),
		arg: objects.ScanArgs{
			Type:   ReleaseMetadata{}.ObjectType(),
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
		if conn, err := r.Query(
			"localnode", "user.list_siblings", map[string]any{
				"out":  "json",
				"zone": astral.ZoneAll,
			}); err != nil {
			for {
				id, _, err := astralRead[*astral.Identity](conn)
				if err != nil {
					return
				}
				a.scan(ctx, *id)
			}
		}
	}()
	return
}

// Read reads an object from the reader using DefaultBlueprints.
func astralRead[O astral.Object](r io.Reader) (o O, n int64, err error) {
	obj, n, err := astral.Read(r)
	if err != nil {
		return
	}
	o, ok := obj.(O)
	if !ok {
		err = fmt.Errorf("invalid object type %s", reflect.TypeOf(obj))
	}
	return
}

type availableAppsScanner struct {
	astrald.Client
	wg  sync.WaitGroup
	log plog.Logger
	rpc rpc.Rpc
	arg objects.ScanArgs
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
	metadata := ReleaseMetadata{}
	readArgs := objects.ReadArgs{
		ID:   id,
		Zone: astral.ZoneAll,
	}
	if err := oc.Fetch(readArgs, &metadata); err != nil {
		a.log.Println(err)
		return
	}

	readArgs.ID = *metadata.ManifestID
	ma := Manifest{}
	if err := oc.Fetch(readArgs, &ma); err != nil {
		a.log.Println(err)
		return
	}
	i := ReleaseInfo{
		Manifest:        ma,
		ReleaseMetadata: metadata,
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
