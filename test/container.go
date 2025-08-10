package test

import (
	"bufio"
	"fmt"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type Container struct {
	id      int
	image   string
	network string
	logfile string
	Args    []any
	io.Writer
	Astrald
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

func (c *Container) Arg(args ...any) *Container {
	cc := *c
	cc.Args = args
	return &cc
}

func (c *Container) Test(run func(t *testing.T), require ...test.Test) test.Test {
	name := fmt.Sprintf("%d-%s", c.id, test.CallerName(2))
	if len(c.Args) == 1 {
		name = fmt.Sprintf("%s%v", name, c.Args[0])
	} else if len(c.Args) > 1 {
		name = fmt.Sprintf("%s%v", name, c.Args)
	}
	return test.New(name, run, require...)
}

func (c *Container) RunContainer() test.Test {
	return c.Test(func(t *testing.T) {
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

func (c *Container) BuildBaseImage() test.Test {
	image := c.image + "-base"
	name := fmt.Sprintf("%s %s", test.CallerName(), image)
	return test.New(name, func(t *testing.T) {
		Command("docker", "build", "-t", image+":latest", "-f", "base.dockerfile", ".").RunT(t)
	})
}

func (c *Container) BuildImage() test.Test {
	//return c.buildImageLocal()
	return c.BuildImageFast()
	//return c.buildImageClean()
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
	return c.Test(func(t *testing.T) {},
		c.StopContainer(),
		c.RemoveContainer(),
	)
}

func (c *Container) StopContainer() test.Test {
	return c.Test(func(t *testing.T) {
		ForceStopContainers(c)
	})
}

func ForceStopContainers(c ...*Container) {
	for _, cc := range c {
		_ = Command("docker", "stop", "-t", "0", cc.Name()).Run()
	}
}

func (c *Container) RemoveContainer() test.Test {
	return c.Test(func(t *testing.T) {
		_ = Command("docker", "rm", c.Name()).Run()
	})
}

func (c *Container) StartLogging() {
	println(fmt.Sprintf(">>> START LOGGING %d", c.id))
	_ = Command("docker", "exec", c.Name(), "sh", "-c", "tail -f "+c.logfile).Run()
	println(fmt.Sprintf(">>> STOP LOGGING %d", c.id))
}

func (c *Container) PrintLog(args ...any) test.Test {
	return c.Arg(args...).Test(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %d", c.id))
		err := Command("docker", "exec", c.Name(), "sh", "-c", "cat "+c.logfile).Run()
		println(fmt.Sprintf(">>> END PRINT LOG %d", c.id))
		assert.NoError(t, err)
	})
}

func (c *Container) ParseLogfile(parsers ...func(log string) bool) test.Test {
	return c.Test(func(t *testing.T) {
		pr, pw := io.Pipe()
		cmd := Command("docker", "exec", c.Name(), "sh", "-c", "tail -n +1 -f "+c.logfile)
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
