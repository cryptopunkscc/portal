package win10

import (
	"context"
	"fmt"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/test/util"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type VirtualMachine struct {
	name   string
	ip     string
	dir    string
	logs   *os.File
	ctx    context.Context
	cancel context.CancelFunc
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{name: "portal-test-win10"}
}

func (vm *VirtualMachine) Name() string {
	return vm.name
}

func (vm *VirtualMachine) Installer() string {
	return `install-portal-to-astral.exe`
}

func (vm *VirtualMachine) Test() test.Test {
	name := fmt.Sprintf("%s-%s", vm.name, test.CallerName(2))
	return test.New(name, func(t *testing.T) {})
}

func (vm *VirtualMachine) PrintLogs() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		file, err := os.ReadFile(vm.logFileName())
		assert.NoError(t, err)
		_, _ = os.Stdout.Write(file)
	})
}

func (vm *VirtualMachine) BuildInstaller() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		root, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		c := util.Command("sh", "-c", "./mage goos windows build:installer")
		c.Dir = root
		c.RunT(t)
	})
}

func (vm *VirtualMachine) init() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		root, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		vm.dir = filepath.Join(root, "test", "kvm", "win10")
		vm.ctx, vm.cancel = context.WithCancel(context.Background())
	})
}

func (vm *VirtualMachine) startVm() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		_ = util.Command("virsh", "start", vm.name).NoStd().Run()
	}).Requires(
		vm.init(),
	)
}

func (vm *VirtualMachine) ParseIp() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		o, err := execCmd("virsh", "domifaddr", vmName).Output()
		assert.NoError(t, err)
		cidr := parseDomifaddr(t, string(o), 0, 3)
		vm.ip = strings.Split(cidr, "/")[0]
		t.Log("parsed IP:", vm.ip)
	}).Requires(
		vm.startVm(),
	)
}

func (vm *VirtualMachine) CopyInstaller() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		root, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		e := "install-portal-to-astral.exe"
		i := filepath.Join(root, "bin", e)
		vm.Scp(i, e).RunT(t)
	}).Requires(
		vm.BuildInstaller(),
		vm.ParseIp(),
	)
}

func (vm *VirtualMachine) RemovePortal() test.Test {
	return vm.Test().Func(func(t *testing.T) {
		_ = vm.Command(vm.Installer(), "-remove").Run()
	})
}

func (vm *VirtualMachine) Start() test.Test {
	return vm.Test().Requires(
		vm.CopyInstaller(),
		vm.RemovePortal(),
	)
}

func (vm *VirtualMachine) Stop() {
	plog.Println("closing portal")
	err := vm.Command("portal", "close").Run()
	if err != nil {
		plog.Println(err)
	}
	time.Sleep(2 * time.Second)
	plog.Println("closing logs")
	err = vm.logs.Close()
	if err != nil {
		plog.Println(err)
	}
	plog.Println("cancelling context")
	vm.cancel()
	plog.Println("awaiting close...")
	time.Sleep(1 * time.Second)
	plog.Println("done")
}

func (vm *VirtualMachine) Command(s ...string) *util.Cmd {
	c := strings.Join(s, " ")
	println("ssh", "user@"+vm.ip, c)
	return util.CommandContext(vm.ctx, "ssh", "user@"+vm.ip, c)
}

func (vm *VirtualMachine) Scp(in, out string) *util.Cmd {
	println("scp", in, "user@"+vm.ip+":"+out)
	return util.Command("scp", in, "user@"+vm.ip+":"+out)
}

func (vm *VirtualMachine) CreateLogFile(t *testing.T) {
	ln := vm.logFileName()
	_ = os.Remove(ln)
	lf, err := os.Create(ln)
	assert.NoError(t, err)
	vm.logs = lf
}

func (vm *VirtualMachine) StartLogging() {
	/*no-op*/
}

func (vm *VirtualMachine) logFileName() string {
	return vm.name + ".log"
}

func (vm *VirtualMachine) StartPortal(t *testing.T) {
	c := vm.Command(`portald`)
	c.Stdout = io.MultiWriter(os.Stdout, vm.logs)
	//err := c.Start()
	//assert.NoError(t, err)
	go func() {
		_ = c.Run()
		plog.Println("portald process done")
	}()
}

func (vm *VirtualMachine) FollowLogFile() *util.Cmd {
	return util.CommandContext(vm.ctx, "sh", "-c", "tail -n +1 -f "+vm.logFileName())
}

func (vm *VirtualMachine) GetBundleName(t *testing.T, dir string, pkg string) string {
	output, err := vm.Command("dir", "build").NoStd().Output()
	assert.NoError(t, err)
	for _, l := range strings.Split(string(output), "\n") {
		c := strings.Split(l, pkg)
		p := strings.Split(c[0], " ")
		if len(c) > 1 {
			return p[len(p)-1] + pkg + c[1]
		}
	}
	//getwd, err := os.Getwd()
	//assert.NoError(t, err)
	//t.Log("getwd:", getwd)
	//entries, err := os.ReadDir(dir)
	//assert.NoError(t, err)
	//for _, f := range entries {
	//	if strings.HasPrefix(f.Name(), pkg) {
	//		return f.Name()
	//	}
	//}
	t.Fatal("bundle not found in:", dir)
	return ""
}

func parseDomifaddr(t *testing.T, s string, r int, c int) string {
	l := strings.Split(strings.TrimSpace(s), "\n")
	assert.Greater(t, len(l), 2)
	f := strings.Fields(l[r+2])
	assert.GreaterOrEqual(t, len(f), c)
	return f[c]
}
