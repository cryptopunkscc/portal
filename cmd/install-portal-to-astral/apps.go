package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/cryptopunkscc/portal/api/env/desktop"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/pkg/config"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func copyApps() (dir string, err error) {
	defer plog.TraceErr(&err)
	dir = filepath.Join(os.TempDir(), "install-portal-to-astral", "apps")
	_ = os.RemoveAll(dir)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	if err = os.CopyFS(dir, apps.Builds); err != nil {
		return
	}
	return
}

func removeDirs() (err error) {
	println(fmt.Sprintf("removing configs..."))
	defer plog.TraceErr(&err)
	c := portal.Config{}
	if err = c.Load(); err != nil {
		if !errors.Is(err, config.ErrNotFound) {
			return // abort when config exist but cannot be loaded for some reason
		}
	}
	if err = c.Build(); err != nil {
		return
	}
	plog.D().Scope("config").Printf("\n%s", c.Yaml())
	for _, s := range c.GetDirs() {
		print(fmt.Sprintf("* removing %s", s))
		err := os.RemoveAll(s)
		print(" [DONE]")
		if err != nil {
			print(fmt.Sprintf(" - %s", err.Error()))
		}
		println()
	}
	return
}
