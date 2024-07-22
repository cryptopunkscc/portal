package main

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	"strings"
	"time"
)

var goWails = GoPortal{
	Tags:    "desktop,wv2runtime.download,production,webkit2_41",
	LdFlags: "-w -s",
	Prepare: func() error { return deps.AptInstallMissing(wailsDeps) },
}

var goWailsDev = GoPortal{
	Tags:    "dev,webkit2_41",
	Prepare: func() error { return deps.AptInstallMissing(wailsDeps) },
}

var goPortalTray = GoPortal{
	Target:  "portal-tray",
	Prepare: func() error { return deps.AptInstallMissing(trayDeps) },
}

var goPortal = GoPortal{}.target("portal")
var goPortalApp = GoPortal{}.target("portal-app")
var goPortalAppWails = goWails.target("portal-app-wails")
var goPortalAppGoja = GoPortal{}.target("portal-app-goja")
var goPortalDev = GoPortal{}.target("portal-dev")
var goPortalDevExec = GoPortal{}.target("portal-dev-exec")
var goPortalDevWails = goWailsDev.target("portal-dev-wails")
var goPortalDevGoja = GoPortal{}.target("portal-dev-goja")
var goPortalDevGo = GoPortal{}.target("portal-dev-go")

func buildPortalInstaller() { GoPortal{Out: "./bin/"}.target("portal-installer").Build() }

type GoPortal struct {
	Cmd     string
	Tags    string
	LdFlags string
	Out     string
	Target  string
	Args    []string
	Prepare func() error
}

func (g GoPortal) Install() {
	g.Cmd = "install"
	if g.prepare() {
		g.run()
	}
}

func (g GoPortal) Build() {
	g.Cmd = "build"
	if g.Out == "" {
		g.Out = "./cmd/portal-installer/bin/"
	}
	g.Args = g.arg("-o", g.Out)
	if g.prepare() {
		g.run()
	}
}

func (g GoPortal) prepare() bool {
	if g.Prepare == nil {
		return true
	}
	if err := g.Prepare(); err != nil {
		log.Printf("cannot %s %s: %v", g.Cmd, g.Target, err)
		return false
	}
	return true
}

func (g GoPortal) target(target string) GoPortal {
	g.Target = target
	return g
}

func (g GoPortal) run() {
	target := fmt.Sprintf("./cmd/%s", g.Target)
	g.Args = append([]string{"go", g.Cmd}, g.Args...)
	g.Args = g.arg("-tags", g.Tags)
	g.Args = g.arg("-ldflags", g.LdFlags)
	g.Args = append(g.Args, target)
	log.Printf("$ %s", strings.Join(g.Args, " "))
	if err := exec.Run(".", g.Args...); err != nil {
		log.Fatalf("%s %s failed: %v", g.Target, g.Cmd, err)
	}
	log.Printf("%s %s succeed.", g.Target, g.Cmd)
}

func (g GoPortal) arg(key string, value string) []string {
	if value != "" {
		return append(g.Args, key, value)
	}
	return g.Args
}

func gpgSignPortalInstaller() {
	time.Sleep(1 * time.Second)
	_ = exec.Run("./bin", "gpg",
		"--sign",
		"--detach-sign",
		"--verbose",
		"--digest-algo", "sha512",
		"./portal-installer",
	)
}
