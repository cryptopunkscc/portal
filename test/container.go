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
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type container struct {
	id      int
	image   string
	network string
	logfile string
	Args    []any
	io.Writer
	astrald
}

func create(amount int, base container) (cs []*container) {
	cs = make([]*container, amount)
	for i := 0; i < amount; i++ {
		cs[i] = base.new()
	}
	return
}

func (c *container) new() *container {
	cc := *c
	c.id++
	return &cc
}

func (c *container) name() string {
	return fmt.Sprintf("%s-%d", c.image, c.id)
}

func (c *container) exec(args ...string) *exec.Cmd {
	return execCmd("docker", append([]string{"exec", "-it", c.name()}, args...)...)
}

func (c *container) execRun(t *testing.T, args ...string) {
	execCmdRun(t, "docker", append([]string{"exec", "-it", c.name()}, args...)...)
}

func (c *container) execRunSh(t *testing.T, args ...string) {
	c.execRun(t, append([]string{"sh", "-c"}, args...)...)
}

func (c *container) execRunShTest(args ...string) test.Test {
	args = append([]string{"sh", "-c"}, args...)
	name := fmt.Sprintf("%s  %s", c.name(), strings.Join(args, " "))
	return test.New(name, func(t *testing.T) {
		c.execRun(t, args...)
	})
}

func (c *container) args(args ...any) *container {
	cc := *c
	cc.Args = args
	return &cc
}

func (c *container) test(run func(t *testing.T), require ...test.Test) test.Test {
	name := fmt.Sprintf("%d-%s", c.id, test.CallerName(2))
	if len(c.Args) == 1 {
		name = fmt.Sprintf("%s%v", name, c.Args[0])
	} else if len(c.Args) > 1 {
		name = fmt.Sprintf("%s%v", name, c.Args)
	}
	return test.New(name, run, require...)
}

func (c *container) runContainer() test.Test {
	return c.test(func(t *testing.T) {
		execCmdRun(t, "docker", "run", "-dit",
			"--rm", // remove container immediately after run
			"--name", c.name(),
			"--network", c.network,
			"-v", "/home/jan/projects/cryptopunks/portal/bin:/portal/bin",
			c.image)
	},
		c.teardown(),
		c.buildImage(),
		c.createNetwork(),
	)
}

func (c *container) buildBaseImage() test.Test {
	image := c.image + "-base"
	name := fmt.Sprintf("%s %s", test.CallerName(), image)
	return test.New(name, func(t *testing.T) {
		execCmdRun(t, "docker", "build", "-t", image+":latest", "-f", "base.dockerfile", ".")
	})
}

func (c *container) buildImage() test.Test {
	//return c.buildImageLocal()
	return c.buildImageFast()
	//return c.buildImageClean()
}

func (c *container) buildImageLocal() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		execCmdRun(t, "docker", "build",
			"--force-rm", "--squash",
			"-f", "local.dockerfile",
			"-t", c.image+":latest", ".")
	},
		c.buildBaseImage(),
		c.removeImage(),
		buildInstaller(),
	)
}

func (c *container) buildImageFast() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		execCmdRun(t, "docker", "build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", projectDir+":/portal",
			"-f", "fast.dockerfile",
			"-t", c.image+":latest", ".")
	},
		c.buildBaseImage(),
		c.removeImage(),
	)
}

func (c *container) buildImageClean() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		cacheDir, err := os.UserCacheDir()
		assert.NoError(t, err)

		projectDir, err := golang.FindProjectRoot()
		assert.NoError(t, err)

		binDir := filepath.Join(projectDir, "bin")
		_ = os.MkdirAll(binDir, 0755)

		execCmdRun(t, "docker", "build",
			"--force-rm", "--squash",
			"-v", build.Default.GOPATH+":/go",
			"-v", cacheDir+":/root/.cache",
			"-v", binDir+":/portal/bin",
			"-f", "clean.dockerfile",
			"-t", c.image+":latest", ".")
	},
		c.buildBaseImage(),
		c.removeImage(),
		packProject(),
	)
}

func (c *container) createNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.network)
	return test.New(name, func(t *testing.T) {
		execCmdRun(t, "docker", "network", "create", c.network)
	},
		c.removeNetwork(),
	)
}

func (c *container) removeNetwork() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.network)
	return test.New(name, func(t *testing.T) {
		_ = execCmd("docker", "network", "rm", "-f", c.network).Run()
	})
}

func (c *container) removeImage() test.Test {
	name := fmt.Sprintf("%s %s", test.CallerName(), c.image)
	return test.New(name, func(t *testing.T) {
		_ = execCmd("docker", "rmi", c.image).Run()
	})
}

func (c *container) teardown() test.Test {
	return c.test(func(t *testing.T) {},
		c.stopContainer(),
		c.removeContainer(),
	)
}

func (c *container) stopContainer() test.Test {
	return c.test(func(t *testing.T) {
		forceStopContainers(c)
	})
}

func forceStopContainers(c ...*container) {
	for _, cc := range c {
		_ = execCmd("docker", "stop", "-t", "0", cc.name()).Run()
	}
}

func (c *container) removeContainer() test.Test {
	return c.test(func(t *testing.T) {
		_ = execCmd("docker", "rm", c.name()).Run()
	})
}

func (c *container) startLogging() {
	println(fmt.Sprintf(">>> START LOGGING %d", c.id))
	_ = execCmd("docker", "exec", c.name(), "sh", "-c", "tail -f "+c.logfile).Run()
	println(fmt.Sprintf(">>> STOP LOGGING %d", c.id))
}

func (c *container) printLog(args ...any) test.Test {
	return c.args(args...).test(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %d", c.id))
		err := execCmd("docker", "exec", c.name(), "sh", "-c", "cat "+c.logfile).Run()
		println(fmt.Sprintf(">>> END PRINT LOG %d", c.id))
		assert.NoError(t, err)
	})
}

func (c *container) parseLogfile(parsers ...func(log string) bool) test.Test {
	return c.test(func(t *testing.T) {
		pr, pw := io.Pipe()
		cmd := execCmd("docker", "exec", c.name(), "sh", "-c", "tail -n +1 -f "+c.logfile)
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
