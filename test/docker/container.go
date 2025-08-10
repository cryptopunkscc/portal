package docker

import (
	"fmt"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/util"
	"github.com/stretchr/testify/assert"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Container struct {
	Id      int
	Image   string
	Network string
	Logfile string
	Args    []any
	RootDir string
	io.Writer
}

func Create(amount int, base Container) (cs []*Container) {
	cs = make([]*Container, amount)
	for i := 0; i < amount; i++ {
		cs[i] = base.new()
	}
	return
}

func (c *Container) Root() string {
	return c.RootDir
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
	return util.Command("docker", "exec", "-it", c.Name(), "sh", "-c", strings.Join(args, " "))
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
		util.Command("docker", "run", "-dit",
			"--rm", // remove Container immediately after run
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
	util.Command("docker", "exec", "-d", c.Name(), "sh", "-c", "portal >> "+c.Logfile+" 2>&1").RunT(t)
}

func (c *Container) CreateLogFile(t *testing.T) {
	c.Command("touch", c.Logfile).RunT(t)
}

func (c *Container) BuildBaseImage() test.Test {
	image := c.Image + "-base"
	name := fmt.Sprintf("%s %s", test.CallerName(), image)
	return test.New(name, func(t *testing.T) {
		util.Command("docker", "build",
			"-f", "docker/base.dockerfile",
			"-t", image+":latest", ".").RunT(t)
	})
}

func (c *Container) BuildImage() test.Test {
	//return c.BuildImageLocal()
	return c.BuildImageFast()
	//return c.BuildImageClean()
}

func (c *Container) BuildImageLocal() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		util.Command("docker", "build",
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

func (c *Container) BuildImageFast() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		util.Command("docker", "build",
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

func (c *Container) BuildImageClean() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		binDir := filepath.Join(projectDir, "bin")
		_ = os.MkdirAll(binDir, 0755)

		util.Command("docker", "build",
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
		util.Command("docker", "network", "create", c.Network).RunT(t)
	},
		c.RemoveNetwork(),
	)
}

func (c *Container) RemoveNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Network)
	return test.New(name, func(t *testing.T) {
		_ = util.Command("docker", "network", "rm", "-f", c.Network).Run()
	})
}

func (c *Container) RemoveImage() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.Image)
	return test.New(name, func(t *testing.T) {
		_ = util.Command("docker", "rmi", c.Image).Run()
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
		ForceStopContainers(c)
	})
}

func ForceStopContainers(c ...*Container) {
	for _, cc := range c {
		_ = util.Command("docker", "stop", "-t", "0", cc.Name()).Run()
	}
}

func (c *Container) RemoveContainer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		_ = util.Command("docker", "rm", c.Name()).Run()
	})
}

func (c *Container) StartLogging() {
	println(fmt.Sprintf(">>> START LOGGING %d", c.Id))
	_ = util.Command("docker", "exec", c.Name(), "sh", "-c", "tail -f "+c.Logfile).Run()
	println(fmt.Sprintf(">>> STOP LOGGING %d", c.Id))
}

func (c *Container) FollowLogFile() *util.Cmd {
	return util.Command("docker", "exec", c.Name(), "sh", "-c", "tail -n +1 -f "+c.Logfile)
}

func (c *Container) PrintLog(args ...any) test.Test {
	return c.Test().Args(args...).Func(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %d", c.Id))
		err := util.Command("docker", "exec", c.Name(), "sh", "-c", "cat "+c.Logfile).Run()
		println(fmt.Sprintf(">>> END PRINT LOG %d", c.Id))
		assert.NoError(t, err)
	})
}
