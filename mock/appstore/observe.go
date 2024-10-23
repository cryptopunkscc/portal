package appstore

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"io"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	plog.Get(ctx).Println("Observing...")
	resolve := apps.Resolver[target.Bundle_]()
	if err = send(resolve, conn, portalAppsDir); err != nil {
		return
	}
	watch, err := fs2.NotifyWatch(ctx, portalAppsDir, 0)
	if err != nil {
		return
	}
	for event := range watch {
		err = send(resolve, conn, event.Name)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			continue
		}
	}

	return
}

func send(
	resolve target.Resolve[target.Bundle_],
	conn rpc.Conn,
	src string,
) (err error) {
	file, err := source.File(src)
	if err != nil {
		return err
	}
	for _, t := range resolve.List(file) {
		if err = conn.Encode(t.Manifest()); err != nil {
			return
		}
	}
	return
}
