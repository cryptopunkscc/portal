package appstore

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"io"
	"log"
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
	targets, err := runner.BundleTargets(src)
	if err != nil {
		return
	}
	for _, target := range targets {
		log.Println("Sending manifest for target", target.Path)
		m := bundle.Manifest{}
		if err := m.LoadFs(target.Files, bundle.PortalJson); err != nil {
			continue
		}
		log.Println("Sending manifest", m)
		if err = conn.Encode(m); err != nil {
			return
		}
	}
	return
}
