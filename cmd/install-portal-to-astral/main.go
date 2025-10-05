package main

import (
	"os"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() {
	cli.Run(handler)
}

var handler = cmd.Handler{
	Name: "install-portal-to-astral",
	Desc: "Astrald & portal environment installer.",
	Func: run,
	Params: cmd.Params{
		{
			Name: "remove",
			Type: "bool",
			Desc: "Remove portal environment.",
		},
		{
			Type: "string",
			Desc: "Optional user name. When specified, the installed node will be assigned to a new user identity associated with the name. Otherwise, the installed node will be ready to claim by existing user.",
		},
	},
}

type Opts struct {
	Remove bool `cli:"remove"`
}

func run(opts Opts, username string) (err error) {
	switch {
	case opts.Remove:
		return remove()
	default:
		return install(username)
	}
}

func install(username string) (err error) {
	if err = installBinaries(); err != nil {
		return
	}

	apps, err := copyApps()
	if err != nil {
		return err
	}
	defer os.RemoveAll(apps)

	setupCmd := []string{"setup", "-apps", apps}
	if len(username) > 0 {
		setupCmd = append(setupCmd, "-user", username)
	}

	if err = portalRun(); err != nil {
		return
	}
	if err = portalRun(setupCmd...); err != nil {
		plog.Println("setup failed:", err)
	}
	if err = portalRun("close"); err != nil {
		return
	}
	return
}

func remove() error {
	if err := removeBinaries(); err != nil {
		return err
	}
	if err := removeDirs(); err != nil {
		return err
	}
	return nil
}
