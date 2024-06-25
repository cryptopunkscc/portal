package appstore

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/pkg/fs"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"io"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	plog.Get(ctx).Println("Observing...")
	err = send(portalAppsDir, conn)
	if err != nil {
		return
	}
	watch, err := fs.NotifyWatch(ctx, portalAppsDir, 0)
	if err != nil {
		return
	}
	for event := range watch {
		err = send(event.Name, conn)
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
	src string,
	conn rpc.Conn,
) (err error) {
	for _, t := range apps.FromPath[target.Bundle](src) {
		if err = conn.Encode(t.Manifest()); err != nil {
			return
		}
	}
	return
}
