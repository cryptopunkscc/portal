package host

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/deprecated/util"
	"github.com/stretchr/testify/assert"
)

func Linux() (os *OS) {
	os = new(OS)
	os.name = "linux"
	os.installer = "./bin/install-portal-to-astral"
	os.ctx, os.cancel = context.WithCancel(context.Background())
	return
}

type OS struct {
	name      string
	installer string
	logs      *os.File
	ctx       context.Context
	cancel    context.CancelFunc
	root      string
	cmd       *util.Cmd
}

func (o *OS) Name() string { return "host-" + o.name }

func (o *OS) Installer() string {
	return o.path(o.installer)
}

func (o *OS) Test() test.Test {
	name := fmt.Sprintf("%s-%s", o.Name(), test.CallerName(2))
	return test.New(name, func(t *testing.T) {})
}

func (o *OS) Start() test.Test {
	return o.Test().Requires(
		BuildInstaller(),
		o.RemovePortal(),
	)
}

func (o *OS) RemovePortal() test.Test {
	return o.Test().Func(func(t *testing.T) {
		_ = o.Command(o.Installer(), "-remove").Run()
	})
}

func (o *OS) PrintLogs() test.Test {
	return o.Test().Func(func(t *testing.T) {
		println(fmt.Sprintf(">>> BEGIN PRINT LOG %s", o.name))
		//_ = o.logs.Close()
		bytes, err := os.ReadFile(o.logFileName())
		assert.NoError(t, err)
		_, _ = os.Stdout.Write(bytes)
		println(fmt.Sprintf(">>> END PRINT LOG %s", o.name))
	})
}

func (o *OS) Stop() {
	if o.cmd == nil {
		return
	}
	var err error
	if o.cmd.ProcessState == nil {
		plog.Println("closing portal")
		time.Sleep(500 * time.Millisecond)
		err = o.Command("portal", "close").Run()
		if err != nil {
			plog.Println(err)
		}
	}
	plog.Println("closing logs")
	_ = o.logs.Close()
	plog.Println("cancelling context")
	o.cancel()
	plog.Println("awaiting close...")
	time.Sleep(100 * time.Millisecond)
	plog.Println("done")
}

func (o *OS) Command(args ...string) *util.Cmd {
	println(strings.Join(args, " "))
	c := util.CommandContext(o.ctx, "sh", "-c", strings.Join(args, " "))
	c.Env = append(c.Env, "ENABLE_PORTAL_APPHOST_LOG=true")
	c.Dir = o.path("test", "host")
	return c
}

func (o *OS) CreateLogFile(t *testing.T) {
	ln := o.logFileName()
	_ = os.Remove(ln)
	lf, err := os.Create(ln)
	assert.NoError(t, err)
	o.logs = lf
}

func (o *OS) logFileName() string {
	return filepath.Join("host", o.Name()+".log")
}

func (o *OS) StartLogging() {}

func (o *OS) StartPortal(t *testing.T) {
	o.cmd = o.Command("portald")
	o.cmd.Stdout = io.MultiWriter(os.Stdout, o.logs)
	o.cmd.Stderr = io.MultiWriter(os.Stderr, o.logs)
	err := o.cmd.Start()
	assert.NoError(t, err)
	go func() {
		_ = o.cmd.Wait()
		plog.Println("portald process done")
	}()
}

func (o *OS) FollowLogFile() (c *util.Cmd) {
	c = o.Command("tail -n +1 -f " + o.logFileName())
	c.Dir = o.path("test")
	return
}

func (o *OS) dir(path ...string) string {
	return o.path("test", "host", filepath.Join(path...))
}

func (o *OS) path(path ...string) string {
	if len(o.root) == 0 {
		root, err := golang.FindProjectRoot()
		if err != nil {
			panic(err)
		}
		o.root = root
	}
	path = append([]string{o.root}, path...)
	return filepath.Join(path...)
}

func (o *OS) GetBundleName(t *testing.T, dir string, pkg string) string {
	getwd, err := os.Getwd()
	plog.Println("getwd", getwd)
	assert.NoError(t, err)
	files, err := os.ReadDir(o.dir(dir))
	assert.NoError(t, err)
	for _, file := range files {
		l := file.Name()
		c := strings.Split(l, pkg)
		p := strings.Split(c[0], " ")
		if len(c) > 1 {
			return p[len(p)-1] + pkg + c[1]
		}
	}
	t.Fatal("bundle not found in:", dir)
	return ""
}

func BuildInstaller() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		cc := util.Command("./mage", "build:installer")
		cc.Dir = "../"
		err := cc.Run()
		assert.NoError(t, err)
	})
}
