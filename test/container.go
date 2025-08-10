package test

import (
	"fmt"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Container struct {
	id      int
	image   string
	network string
	logfile string
	Args    []any
	root    string
	io.Writer
	Astrald
}

func (c *Container) Root() string {
	return c.root
}

func Create(amount int, base Container) (cs []*Container) {
	cs = make([]*Container, amount)
	for i := 0; i < amount; i++ {
		cs[i] = base.new()
	}
	return
}

func (c *Container) new() *Container {
	cc := *c
	c.id++
	return &cc
}

func (c *Container) Name() string {
	return fmt.Sprintf("%s-%d", c.image, c.id)
}

func (c *Container) Command(args ...string) *Cmd {
	return Command("docker", "exec", "-it", c.Name(), "sh", "-c", strings.Join(args, " "))
}

func (c *Container) GetBundleName(t *testing.T, dir string, pkg string) string {
	b := c.Command("ls", dir, "| grep", pkg).OutputT(t)
	return strings.TrimSpace(string(b))
}

func (c *Container) Test() test.Test {
	name := fmt.Sprintf("%d-%s", c.id, test.CallerName(2))
	return test.New(name, func(t *testing.T) {})
}

func (c *Container) Start() test.Test {
	return c.Test().Func(func(t *testing.T) {
		Command("docker", "run", "-dit",
			"--rm", // remove Container immediately after run
			"--name", c.Name(),
			"--network", c.network,
			"-v", "/home/jan/projects/cryptopunks/portal/bin:/portal/bin",
			c.image,
		).RunT(t)
	},
		c.Teardown(),
		c.BuildImage(),
		c.CreateNetwork(),
	)
}

func (c *Container) StartPortal(t *testing.T) {
	Command("docker", "exec", "-d", c.Name(), "sh", "-c", "portal >> "+c.logfile+" 2>&1").RunT(t)
}

func (c *Container) CreateLogFile(t *testing.T) {
	c.Command("touch", c.logfile).RunT(t)
}

func (c *Container) BuildBaseImage() test.Test {
	image := c.image + "-base"
	name := fmt.Sprintf("%s %s", test.CallerName(), image)
	return test.New(name, func(t *testing.T) {
		Command("docker", "build", "-t", image+":latest", "-f", "base.dockerfile", ".").RunT(t)
	})
}

func (c *Container) BuildImage() test.Test {
	//return c.BuildImageLocal()
	return c.BuildImageFast()
	//return c.BuildImageClean()
}

func (c *Container) BuildImageLocal() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		Command("docker", "build",
			"--force-rm", "--squash",
			"-f", "local.dockerfile",
			"-t", c.image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
		BuildInstaller(),
	)
}

func (c *Container) BuildImageFast() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		Command("docker", "build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", projectDir+":/portal",
			"-f", "fast.dockerfile",
			"-t", c.image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
	)
}

func (c *Container) BuildImageClean() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		binDir := filepath.Join(projectDir, "bin")
		_ = os.MkdirAll(binDir, 0755)

		Command("docker", "build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", binDir+":/portal/bin",
			"-f", "clean.dockerfile",
			"-t", c.image+":latest", ".",
		).RunT(t)
	},
		c.BuildBaseImage(),
		c.RemoveImage(),
		PackProject(),
	)
}

func (c *Container) CreateNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.network)
	return test.New(name, func(t *testing.T) {
		Command("docker", "network", "create", c.network).RunT(t)
	},
		c.RemoveNetwork(),
	)
}

func (c *Container) RemoveNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.network)
	return test.New(name, func(t *testing.T) {
		_ = Command("docker", "network", "rm", "-f", c.network).Run()
	})
}

func (c *Container) RemoveImage() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		_ = Command("docker", "rmi", c.image).Run()
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
		_ = Command("docker", "stop", "-t", "0", cc.Name()).Run()
	}
}

func (c *Container) RemoveContainer() test.Test {
	return c.Test().Func(func(t *testing.T) {
		_ = Command("docker", "rm", c.Name()).Run()
	})
}

func (c *Container) StartLogging() {
	println(fmt.Sprintf(">>> START LOGGING %d", c.id))
	_ = Command("docker", "exec", c.Name(), "sh", "-c", "tail -f "+c.logfile).Run()
	println(fmt.Sprintf(">>> STOP LOGGING %d", c.id))
}

func (c *Container) FollowLogFile() *Cmd {
	return Command("docker", "exec", c.Name(), "sh", "-c", "tail -n +1 -f "+c.logfile)
}

func (c *Container) PrintLog(args ...any) test.Test {
	return c.Test().Args(args...).Func(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %d", c.id))
		err := Command("docker", "exec", c.Name(), "sh", "-c", "cat "+c.logfile).Run()
		println(fmt.Sprintf(">>> END PRINT LOG %d", c.id))
		assert.NoError(t, err)
	})
}
