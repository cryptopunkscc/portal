package test

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"path/filepath"
	"testing"
	"time"
)

func (c *container) printInstallHelp() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "./portal-installer h")
	},
		c.runContainer(),
	)
}

func (c *container) installFirstPortal() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "./portal-installer first test_user")
	},
		c.runContainer(),
	)
}

func (c *container) installNextPortal() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "./portal-installer next")
	},
		c.runContainer(),
	)
}

func (c *container) portalStart() test.Test {
	installPortal := c.installFirstPortal()
	if c.id > 0 {
		installPortal = c.installNextPortal()
	}
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "touch "+c.logfile)
		go c.startLogging()

		execCmdRun(t, "docker", "exec", "-d", c.name(), "sh", "-c", "portal >> "+c.logfile+" 2>&1")
		time.Sleep(1 * time.Second)
	},
		installPortal,
	)
}

func (c *container) portalStartAwait() test.Test {
	return c.test(func(t *testing.T) {},
		c.portalStart(),
		c.parseLogfile(
			c.parseNodeInfo,
		),
	)
}

func (c *container) portalClose() test.Test {
	return c.test(func(t *testing.T) {
		time.Sleep(1000 * time.Millisecond)
		c.execRunSh(t, "portal close")
	},
		c.portalStart(),
	)
}

func (c *container) portalHelp() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "portal h")
	},
		c.portalStart(),
	)
}

func (c *container) createUser() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "portal user create test_user")
	},
		c.portalStart(),
	)
}

func (c *container) userInfo() test.Test {
	return c.test(func(t *testing.T) {
		c.execRunSh(t, "portal user info")
	},
		c.portalStart(),
	)
}

func (c *container) userClaim(c2 *container) test.Test {
	return c.args(c2.name()).test(func(t *testing.T) {
		c.execRun(t, "portal", "user", "claim", c2.identity)
	},
		c.portalStart(),
		c2.portalStartAwait(),
	)
}

func (c *container) listTemplates(runner string) test.Test {
	return c.args(runner).test(func(t *testing.T) {
		c.execRun(t, "portal", runner, "templates")
	},
		c.portalStart(),
	)
}

type projectOpts struct {
	runner   string
	template string
	name     string
}

func (o projectOpts) Name() string {
	if o.name != "" {
		return o.name
	}
	return o.template
}

func (c *container) newProject(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		c.execRun(t, "portal", opts.runner, "new", "-t", opts.template, opts.Name())
	},
		c.portalStart(),
	)
}

func (c *container) buildProject(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		c.execRunSh(t, "ls -lah")
		c.execRun(t, "portal", "build", opts.Name(), "pack", ".")
		c.execRunSh(t, "ls -lah")
		c.execRunSh(t, "ls -lah ./build")
		time.Sleep(1 * time.Second)
	},
		c.newProject(opts),
	)
}

func (c *container) publishProject(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		e := c.exec("sh", "-c", "ls ./build | grep "+opts.Name())
		e.Stdout = nil
		b, err := e.Output()
		test.AssertErr(t, err)
		p := filepath.Join("./build", string(b))
		t.Log("publishing:", p)
		c.execRun(t, "portal", "app", "publish", p)
	},
		c.buildProject(opts),
	)
}

func (c *container) listAvailableApps(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		c.execRunSh(t, "portal app available")
	},
		c.buildProject(opts),
	)
}

func (c *container) installAvailableApp(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		c.execRunSh(t, "portal app install my.app."+opts.Name())
	})
}

func (c *container) runApp(opts projectOpts) test.Test {
	return c.args(opts).test(func(t *testing.T) {
		c.execRunSh(t, "portal "+opts.Name())
	})
}
