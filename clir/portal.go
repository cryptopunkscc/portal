package clir

import (
	"context"
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/leaanthony/clir"
	"log"
)

func RunPortal(
	ctx context.Context,
	bindings runtime.New,
	open Open,
	list List,
	install Install,
	uninstall Uninstall,
	version Version,
) error {
	c := newCli(ctx, bindings)
	c.clir = clir.NewCli(portal.Name, portal.Description, version())
	c.Open(open)
	c.List(list)
	c.Install(install)
	c.Uninstall(uninstall)
	c.Version(version)
	return c.clir.Run()
}

type Open func(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error)

type OpenFlags struct {
	Src    string `pos:"1" default:""`
	Attach bool   `name:"attach" description:"Attach runner to the current process instead of dispatching execution to portal service. If given path contains multiple portals each will be run as a child process."`
}

func (c cli) Open(handle Open) {
	flags := &OpenFlags{}
	f := func() error {
		return handle(c.ctx, c.bindings, flags.Src, flags.Attach)
	}
	cmd := c.clir.NewSubCommand("o", "Open app from a given source. The source can be a app name, package name, app bundle or app dir.")
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}

type List func() []target.App

func (c cli) List(handle List) {
	cmd := c.clir.NewSubCommand("l", "List installed apps.")
	cmd.Action(func() (_ error) {
		for i, app := range handle() {
			log.Println(i, app.Manifest())
		}
		return
	})
	return
}

type Install func(string) error

func (c cli) Install(handle Install) {
	flags := struct {
		Path string `pos:"1" default:""`
	}{}
	cmd := c.clir.NewSubCommand("i", "Install app from a given portal app bundle path.")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		return handle(flags.Path)
	})
	return
}

type Uninstall func(string) error

func (c cli) Uninstall(handle Uninstall) {
	flags := struct {
		Id string `pos:"1" default:""`
	}{}
	cmd := c.clir.NewSubCommand("u", "Uninstall app by given id. The id can by app name or package name")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		return handle(flags.Id)
	})
	return
}
