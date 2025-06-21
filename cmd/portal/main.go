package main

import (
	_ "github.com/cryptopunkscc/portal/api/env/desktop"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
)

func init() {
	plog.Default = plog.Default.Scope("cli")
}

func main() {
	a := Application{}
	cli.Run(a.Handler())
}

type Application struct {
	Config  portal.Config
	Apphost apphost.Adapter
}
