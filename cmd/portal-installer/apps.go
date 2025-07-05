package main

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"path/filepath"
)

func installApps() (err error) {
	defer plog.TraceErr(&err)
	dir := filepath.Join(os.TempDir(), "portal-installer", "apps")
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	defer os.RemoveAll(dir)
	if err = os.CopyFS(dir, apps.Builds); err != nil {
		return
	}
	if err = portalRun("app", "install", dir); err != nil {
		return
	}
	return
}
