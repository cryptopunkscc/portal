package test

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"path/filepath"
	"testing"
	"time"
)

func (c *Container) PrintInstallHelp() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("/portal/bin/install-portal-to-astral h").RunT(t)
	},
		c.RunContainer(),
	)
}

func (c *Container) InstallFirstPortal() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("/portal/bin/install-portal-to-astral test_user").RunT(t)
	},
		c.RunContainer(),
	)
}

func (c *Container) InstallNextPortal() test.Test {
	return c.Test(func(t *testing.T) {
		c.Command("/portal/bin/install-portal-to-astral").RunT(t)
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
		c.Command("touch " + c.logfile).RunT(t)
		go c.StartLogging()

		Command("docker", "exec", "-d", c.Name(), "sh", "-c", "portal >> "+c.logfile+" 2>&1").RunT(t)
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
		c.Command("ls -lah").RunT(t)
		c.Command("portal", "build", "-p", "-o", ".", opts.Name()).RunT(t)
		c.Command("ls -lah").RunT(t)
		c.Command("ls -lah ./build").RunT(t)
		time.Sleep(1 * time.Second)
	},
		c.NewProject(opts),
	)
}

func (c *Container) PublishProject(opts ProjectOpts) test.Test {
	return c.Arg(opts).Test(func(t *testing.T) {
		e := c.Command("ls ./build | grep " + opts.Name())
		e.Stdout = nil
		b, err := e.Output()
		test.AssertErr(t, err)
		p := filepath.Join("./build", string(b))
		t.Log("publishing:", p)
		c.Command("portal", "app", "publish", p).RunT(t)
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
