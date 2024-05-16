package clir

import (
	"context"
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/leaanthony/clir"
	"log"
	"os"
	"reflect"
	"strings"
)

func RunPortalDev(
	ctx context.Context,
	bindings runtime.New,
	dev Dev,
	templates Templates,
	create Create,
	build Build,
	version Version,
) error {
	c := newCli(ctx, bindings)
	c.clir = clir.NewCli(portal.NameDev, portal.DescriptionDev, version())
	c.Dev(dev)
	c.Create(templates, create)
	c.Build(build)
	c.Version(version)
	c.Apps()
	return c.clir.Run()
}

type Dev func(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error)

func (c cli) Dev(handle Dev) {
	flags := &struct {
		Src    string `pos:"1" default:"" description:"Application source. The source can be a app name, package name, app bundle or app dir."`
		Attach bool   `name:"attach" description:"Attach runner to the current process instead of dispatching execution to portal service. If given path contains multiple portals each will be run as a child process."`
	}{}
	f := func() error {
		return handle(c.ctx, c.bindings, flags.Src, flags.Attach)
	}
	cmd := c.clir.NewSubCommand("o", "Open project or app from a given source in development environment.")
	cmd.AddFlags(flags)
	cmd.Action(f)
	c.clir.DefaultCommand(cmd)
	c.clir.AddFlags(flags)
	c.clir.Action(f)
	return
}

type Templates func() error

type Create func(
	projectName string,
	targetDir string,
	templates []string,
	force bool,
) (err error)

func (c cli) Create(
	templates Templates,
	create Create,
) {
	emptyFlags := struct {
		Dir      string `pos:"1" description:"Project directory"`
		Name     string `name:"n" description:"Name of project"`
		Template string `name:"t" description:"Name of built-in template to use, path to template or template url"`
		Force    bool   `name:"f" description:"Force recreate project"`
		List     bool   `name:"l" description:"List available templates"`
	}{}
	flags := emptyFlags
	cmd := c.clir.NewSubCommand("c", "Create new project from template.")
	cmd.AddFlags(&flags)
	cmd.Action(func() error {
		switch {
		case flags == emptyFlags || flags.List:
			return templates()
		default:
			temps := strings.Split(flags.Template, " ")
			return create(flags.Name, flags.Dir, temps, flags.Force)
		}
	})
	return
}

type Build func(string) error

func (c cli) Build(handle Build) {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("b", "Build project and generate portal app bundle.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		return handle(flags.Path)
	})
	return
}

func (c cli) Apps() {
	flags := struct {
		Path string `pos:"1" default:"."`
	}{}
	cmd := c.clir.NewSubCommand("t", "Print all targets in given directory.")
	cmd.AddFlags(&flags)
	cmd.Action(func() (err error) {
		wd, _ := os.Getwd()
		for source := range project.FindInPath[target.Source](flags.Path) {
			log.Println(reflect.TypeOf(source), "\t", strings.TrimPrefix(source.Abs(), wd+"/"))
		}
		return
	})
	return
}
