package test

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func (c *Container) installPortalToAstral(args ...string) *Cmd {
	args = slices.Insert(args, 0, filepath.Join(c.root, "portal", "bin", "install-portal-to-astral"))
	return c.Command(args...)
}

func (c *Container) PrintInstallHelp() test.Test {
	return c.Test(func(t *testing.T) {
		c.installPortalToAstral("h").RunT(t)
	},
		c.RunContainer(),
	)
}

func (c *Container) InstallFirstPortal() test.Test {
	return c.Test(func(t *testing.T) {
		c.installPortalToAstral("test_user").RunT(t)
	},
		c.RunContainer(),
	)
}

func (c *Container) InstallNextPortal() test.Test {
	return c.Test(func(t *testing.T) {
		c.installPortalToAstral().RunT(t)
	},
		c.RunContainer(),
	)
}

func (c *Container) PortalStart() test.Test {
	installPortal := c.InstallFirstPortal()
	if c.id > 0 {
		installPortal = c.InstallNextPortal()
	}
	return c.Test(func(t *testing.T) {
		c.CreateLogFile(t)
		go c.StartLogging()

		c.StartPortal(t)
		time.Sleep(1 * time.Second)
	},
		installPortal,
	)
}

func (c *Container) PortalStartAwait() test.Test {
	return c.Test(func(t *testing.T) {},
		c.PortalStart(),
		c.ParseLogfile(
			c.ParseNodeInfo,
		),
	)
}

func (c *Container) PortalClose() test.Test {
	return c.Test(func(t *testing.T) {
		time.Sleep(1000 * time.Millisecond)
		c.Command("portal close").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Container) PortalHelp() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("portal h").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Container) CreateUser() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("portal user create test_user").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Container) UserInfo() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("portal user info").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Container) UserClaim(c2 *Container) test.Test {
	return c.Arg(c2.Name()).Test(func(t *testing.T) {
		c.Command("portal", "user", "claim", c2.identity).RunT(t)
	},
		c.PortalStart(),
		c2.PortalStartAwait(),
	)
}

func (c *Container) ListTemplates(runner string) test.Test {
	return c.Arg(runner).Test(func(t *testing.T) {
		c.Command("portal", runner, "templates").RunT(t)
	},
		c.PortalStart(),
	)
}

type ProjectOpts struct {
	runner   string
	template string
	name     string
}

func (o ProjectOpts) Name() string {
	if o.name != "" {
		return o.name
	}
	return o.template
}

func (c *Container) NewProject(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		c.Command("portal", opts.runner, "new", "-t", opts.template, opts.Name()).RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Container) BuildProject(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		c.Command("portal", "build", "-p", "-o", ".", opts.Name()).RunT(t)
		time.Sleep(1 * time.Second)
	},
		c.NewProject(opts),
	)
}

func (c *Container) PublishProject(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		name := c.GetBundleName(t, "build", opts.Name())
		path := filepath.Join("build", name)
		t.Log("publishing:", path)
		c.Command("portal", "app", "publish", path).RunT(t)
	},
		c.BuildProject(opts),
	)
}

func (c *Container) ListAvailableApps(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		c.Command("portal app available").RunT(t)
	},
		c.BuildProject(opts),
	)
}

func (c *Container) InstallAvailableApp(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		c.Command("portal app install my.app." + opts.Name()).RunT(t)
	})
}

func (c *Container) RunApp(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		c.Command("portal " + opts.Name()).RunT(t)
	})
}
