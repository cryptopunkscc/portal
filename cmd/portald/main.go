package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"path/filepath"
)

var application = &Application[Portal_]{}

func init() {
	application.ExtraTokens = []string{"portal"}
}

func main() {
	ctx := context.Background()
	log := plog.New().D().Scope("portald").Set(&ctx)
	go singal.OnShutdown(log, application.Stop)

	c := application.commands()
	cmd.InjectHelp(&c)
	err := cli.New(c).Run(ctx)
	if err != nil {
		log.E().Println("finished with error:", err)
	}
}

func init() {
	env.AstraldHome.Default(defaultAstraldHome)
	env.PortaldHome.Default(defaultPortalHome)
}

func defaultAstraldHome() string { return filepath.Join(userConfigDir(), "astrald") }
