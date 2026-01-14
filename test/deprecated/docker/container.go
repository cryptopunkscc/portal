package docker

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/deprecated/util"
	"github.com/stretchr/testify/assert"
)

type Container struct {
	Id            int
	Bin           string
	Build         func(*Container) test.Test
	Image         string
	Network       string
	Logfile       string
	Args          []any
	InstallerPath string
	io.Writer
}

func (c *Container) New(id int) *Container {
	cc := *c
	cc.Id = id
	return &cc
}

func (c *Container) Installer() string {
	return c.InstallerPath
}

func (c *Container) new() *Container {
	cc := *c
	c.Id++
	return &cc
}

func (c *Container) Name() string {
	return fmt.Sprintf("%s-%d", c.Image, c.Id)
}

func (c *Container) Command(args ...string) *util.Cmd {
	return util.Command(c.Bin, "exec", "-t", c.Name(), "sh", "-c", strings.Join(args, " "))
}

func (c *Container) GetBundleName(t *testing.T, dir string, pkg string) string {
	b := c.Command("ls", dir, "| grep", pkg).OutputT(t)
	return strings.TrimSpace(string(b))
}

func (c *Container) Test() test.Test {
	name := fmt.Sprintf("%d-%s", c.Id, test.CallerName(2))
	return test.New(name, func(t *testing.T) {})
}

func (c *Container) Start() test.Test {
	return c.Test().Func(func(t *testing.T) {
		util.Command(c.Bin, "run", "-dit",
			//"--rm", // remove Container immediately after run
			"--name", c.Name(),
			"--network", c.Network,
			"-v", "/home/jan/projects/cryptopunks/portal/bin:/portal/bin",
			c.Image,
		).RunT(t)
	},
		c.Teardown(),
		c.BuildImage(),
		c.CreateNetwork(),
	)
}

func (c *Container) StartPortal(t *testing.T) {
	util.Command(c.Bin, "exec", "-d", c.Name(), "sh", "-c", "portal >> "+c.Logfile+" 2>&1").RunT(t)
}

func (c *Container) CreateLogFile(t *testing.T) {
	c.Command("touch", c.Logfile).RunT(t)
}

func (c *Container) BuildBaseImage() test.Test {
	image := c.Image + "-base"
	name := fmt.Sprintf("%s %s", test.CallerName(), image)
	return test.New(name, func(t *testing.T) {
		util.Command(c.Bin, "build",
			"-f", "docker/base.dockerfile",
			"-t", image+":latest", ".").RunT(t)
	})
}

func (c *Container) BuildImage() test.Test {
	if c.Build == nil {
		return BuildImageLocal(c)
	}
	return c.Build(c)
}

func BuildImageLocal(c *Container) test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		util.Command(c.Bin,
			"buildx",
			"build",
			"--force-rm", "--squash",
			"-f", "docker/local.dockerfile",
			"-t", c.Image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
		BuildInstaller(),
	)
}

func BuildInstaller() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		cc := util.Command("./mage", "build:installer")
		cc.Dir = "../"
		err := cc.Run()
		assert.NoError(t, err)
	})
}

func BuildImageFast(c *Container) test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		util.Command(c.Bin,
			"buildx",
			"build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", projectDir+":/portal",
			"-f", "docker/fast.dockerfile",
			"-t", c.Image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
	)
}

func BuildImageClean(c *Container) test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		binDir := filepath.Join(projectDir, "bin")
		_ = os.MkdirAll(binDir, 0755)

		util.Command(c.Bin,
			"buildx",
			"build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", binDir+":/portal/bin",
			"-f", "docker/clean.dockerfile",
			"-t", c.Image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
		PackProject(),
	)
}

func PackProject() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		root, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		c := util.Command("sh", "-c", "git ls-files -co --exclude-standard -z | tar -cf ./test/docker/sources.tar --exclude=./test/docker/sources.tar --null -T -")
		c.Dir = root
		err = c.Run()
		assert.NoError(t, err)
	})
}

func (c *Container) CreateNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Network)
	return test.New(name, func(t *testing.T) {
		util.Command(c.Bin, "network", "create", c.Network).RunT(t)
	},
		c.RemoveNetwork(),
	)
}

func (c *Container) RemoveNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Network)
	return test.New(name, func(t *testing.T) {
		_ = util.Command(c.Bin, "network", "rm", "-f", c.Network).Run()
	})
}

func (c *Container) RemoveImage() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		_ = util.Command(c.Bin, "rmi", c.Image).Run()
	})
}

func (c *Container) Teardown() test.Test {
	return c.Test().Func(func(t *testing.T) {},
		c.StopContainer(),
		c.RemoveContainer(),
	)
}

func (c *Container) StopContainer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		c.Stop()
	})
}

func (c *Container) Stop() {
	_ = util.Command(c.Bin, "stop", "-t", "0", c.Name()).Run()
}

func (c *Container) RemoveContainer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		_ = util.Command(c.Bin, "rm", c.Name()).Run()
	})
}

func (c *Container) StartLogging() {
	println(fmt.Sprintf(">>> START LOGGING %d", c.Id))
	_ = util.Command(c.Bin, "exec", c.Name(), "sh", "-c", "tail -f "+c.Logfile).Run()
	println(fmt.Sprintf(">>> STOP LOGGING %d", c.Id))
}

func (c *Container) FollowLogFile() *util.Cmd {
	return util.Command(c.Bin, "exec", c.Name(), "sh", "-c", "tail -n +1 -f "+c.Logfile)
}

func (c *Container) PrintLogs() test.Test {
	return c.Test().Func(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %d", c.Id))
		err := util.Command(c.Bin, "exec", c.Name(), "sh", "-c", "cat "+c.Logfile).Run()
		println(fmt.Sprintf(">>> END PRINT LOG %d", c.Id))
		assert.NoError(t, err)
	})
}
