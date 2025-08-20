package test

import (
	"bufio"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/util"
	"github.com/stretchr/testify/assert"
	"io"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

type Device interface {
	Name() string
	Installer() string
	Test() test.Test
	Start() test.Test
	PrintLogs() test.Test
	Stop()
	Command(...string) *util.Cmd

	CreateLogFile(t *testing.T)
	StartLogging()
	StartPortal(t *testing.T)
	FollowLogFile() *util.Cmd
	GetBundleName(t *testing.T, dir string, pkg string) string
}

type Cases struct {
	id int
	Device
	Astrald
}

func (c *Cases) installPortalToAstral(args ...string) *util.Cmd {
	args = slices.Insert(args, 0, c.Installer())
	return c.Command(args...)
}

func (c *Cases) PrintInstallHelp() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.installPortalToAstral("h").RunT(t)
	},
		c.Start(),
	)
}

func (c *Cases) InstallFirstPortal() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.installPortalToAstral("test_user").RunT(t)
		//time.Sleep(200 * time.Millisecond) // windows experimental
	}).Requires(
		c.Start(),
	)
}

func (c *Cases) InstallNextPortal() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.installPortalToAstral().RunT(t)
	}).Requires(
		c.Start(),
	)
}

func (c *Cases) PortalStart() test.Test {
	installPortal := c.InstallFirstPortal()
	if c.id > 0 {
		installPortal = c.InstallNextPortal()
	}
	return c.Test().Func(func(t *testing.T) {
		c.CreateLogFile(t)
		go c.StartLogging()

		c.Device.StartPortal(t)
		time.Sleep(1 * time.Second)
	},
		installPortal,
	)
}

func (c *Cases) PortalStartAwait() test.Test {
	return c.Test().Requires(
		c.PortalStart(),
		c.ParseLogfile(
			c.ParseNodeInfo,
		),
	)
}

func (c *Cases) ParseLogfile(parsers ...func(log string) bool) test.Test {
	return c.Test().Func(func(t *testing.T) {
		pr, pw := io.Pipe()
		cmd := c.FollowLogFile()
		cmd.Stdout = pw
		cmd.Stderr = pw
		if err := cmd.Start(); !assert.NoError(t, err) {
			return
		}

		s := bufio.NewScanner(pr)
		defer func() {
			_ = cmd.Process.Kill()
			_ = pr.Close()
			_ = pw.Close()
		}()
		for s.Scan() {
			l := s.Text()
			for i, parser := range parsers {
				if parser(l) {
					parsers = slices.Delete(parsers, i, i+1)
					break
				}
			}
			if len(parsers) == 0 {
				break
			}
		}
	})
}

func (c *Cases) PortalClose() test.Test {
	return c.Test().Func(func(t *testing.T) {
		time.Sleep(1000 * time.Millisecond)
		c.Command("portal", "close").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Cases) PortalHelp() test.Test {
	return c.Test().Func(func(t *testing.T) {
		//time.Sleep(1 * time.Second) // windows experimental
		c.Command("portal", "h").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Cases) CreateUser() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.Command("portal", "user", "create", "test_user").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Cases) UserInfo() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.Command("portal", "user", "info").RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Cases) UserClaim(c2 *Cases) test.Test {
	return c.Test().Args(c2.Name()).Func(func(t *testing.T) {
		//time.Sleep(3 * time.Second) // windows experimental
		c.Command("portal", "user", "claim", c2.identity).RunT(t)
	},
		c.PortalStart(),
		c2.PortalStartAwait(),
	)
}

func (c *Cases) ListTemplates(runner string) test.Test {
	return c.Test().Args(runner).Func(func(t *testing.T) {
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

func (c *Cases) NewProject(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		c.Command("portal", opts.runner, "new", "-t", opts.template, opts.Name()).RunT(t)
	},
		c.PortalStart(),
	)
}

func (c *Cases) BuildProject(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		c.Command("portal", "build", "-p", "-o", ".", opts.Name()).RunT(t)
		time.Sleep(1 * time.Second)
	},
		c.NewProject(opts),
	)
}

func (c *Cases) PublishProject(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		name := c.GetBundleName(t, "build", opts.Name())
		path := filepath.Join("build", name)
		t.Log("publishing:", path)
		c.Command("portal", "app", "publish", path).RunT(t)
	},
		c.BuildProject(opts),
	)
}

func (c *Cases) ListAvailableApps(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		c.Command("portal", "app", "available").RunT(t)
	},
		c.BuildProject(opts),
	)
}

func (c *Cases) InstallAvailableApp(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		//time.Sleep(1 * time.Second) // windows experimental
		c.Command("portal", "app", "install", "my.app."+opts.Name()).RunT(t)
	})
}

func (c *Cases) RunApp(opts ProjectOpts) test.Test {
	return c.Test().Args(opts).Func(func(t *testing.T) {
		c.Command("portal", opts.Name()).RunT(t)
	})
}
