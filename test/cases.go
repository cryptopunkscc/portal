package test

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"testing"
	"time"
)

func (c *container) printInstallHelp() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "./portal-installer h")
	},
		c.runContainer(),
	)
}

func (c *container) installFirstPortal() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "./portal-installer first test_user")
	},
		c.runContainer(),
	)
}

func (c *container) installNextPortal() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "./portal-installer next")
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
		c.dockerExecSh(t, "touch "+c.logfile)
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
		c.dockerExecSh(t, "portal close")
	},
		c.portalStart(),
	)
}

func (c *container) portalHelp() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "portal h")
	},
		c.portalStart(),
	)
}

func (c *container) createUser() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "portal user create test_user")
	},
		c.portalStart(),
	)
}

func (c *container) userInfo() test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExecSh(t, "portal user info")
	},
		c.portalStart(),
	)
}

func (c *container) claim(c2 *container) test.Test {
	return c.test(func(t *testing.T) {
		c.dockerExec(t, "portal", "user", "claim", c2.identity)
	},
		c.portalStart(),
		c2.portalStartAwait(),
	)
}
