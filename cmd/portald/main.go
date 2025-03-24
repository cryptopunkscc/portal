package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/exec"
	"path/filepath"
)

var application = &Application[Portal_]{}

func main() {
	ctx := context.Background()
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(log, application.Stop)

	env.PortaldBin.MkdirAll()
	application.NodeDir = mem.NewVar(env.AstraldHome.MkdirAll())
	application.AppsDir = mem.NewVar(env.PortaldApps.MkdirAll())
	application.TokensDir = mem.NewVar(env.PortaldTokens.MkdirAll())
	application.Astrald = &exec.Astrald{NodeRoot: application.NodeDir}
	application.CreateTokens = []string{"portal"}

	handler := application.handler()
	cmd.InjectHelp(&handler)
	err := cli.New(handler).Run(ctx)
	if err != nil {
		log.E().Println("finished with error:", err)
	}
}

func init() {
	env.AstraldHome.Default(defaultAstraldHome)
	env.AstraldDb.Default(defaultAstraldHome)
	env.PortaldTokens.Default(defaultTokensDir)
	env.PortaldApps.Default(defaultAppsDir)
	env.PortaldBin.Default(defaultBinDir)
}

func defaultAstraldHome() string { return filepath.Join(userConfigDir(), "astrald") }
func defaultTokensDir() string   { return filepath.Join(defaultPortalDir(), "tokens") }
func defaultBinDir() string      { return filepath.Join(defaultPortalDir(), "bin") }
