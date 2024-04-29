package appstore

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"io"
	"log"
	"os"
)

func Observe(ctx context.Context, conn rpc.Conn) (err error) {
	log.Println("Observing...")

	err = send(portalAppsDir, conn)
	if err != nil {
		return
	}
	watch, err := fs.NotifyWatch(ctx, portalAppsDir)
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
	for target := range project.BundleTargets(os.DirFS(src), ".") {
		m := bundle.Manifest{}
		if err := m.LoadFs(target.Files(), bundle.PortalJson); err != nil {
			continue
		}
		if err = conn.Encode(m); err != nil {
			return
		}
	}
	return
}
