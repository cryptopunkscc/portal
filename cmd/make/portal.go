package main

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	"strings"
)

var goWails = GoPortal{
	Tags:    "desktop,wv2runtime.download,production,webkit2_41",
	LdFlags: []string{"-w -s"},
	Prepare: func() error { return deps.AptInstallMissing(wailsDeps) },
}

var goWailsDev = GoPortal{
	Tags:    "dev,webkit2_41",
	Prepare: func() error { return deps.AptInstallMissing(wailsDeps) },
}

var goPortal = GoPortal{}.target("portal")
var goPortalApp = GoPortal{}.target("portal-app")
var goPortalApps = GoPortal{}.target("portal-apps")
var goPortalAppWails = goWails.target("portal-app-wails")
var goPortalAppGoja = GoPortal{}.target("portal-app-goja")
var goPortalNew = GoPortal{}.target("portal-new")
var goPortalList = GoPortal{}.target("portal-list")
var goPortalBuild = GoPortal{}.target("portal-build")
var goPortalDev = GoPortal{}.target("portal-dev")
var goPortalDevExec = GoPortal{}.target("portal-dev-exec")
var goPortalDevWails = goWailsDev.target("portal-dev-wails")
var goPortalDevGoja = GoPortal{}.target("portal-dev-goja")
var goPortalDevGo = GoPortal{}.target("portal-dev-go")

func buildPortalInstaller() { GoPortal{Out: "./bin/"}.target("portal-installer").Build() }

type GoPortal struct {
	Cmd     string
	Tags    string
	LdFlags []string
	Out     string
	Target  string
	Args    []string
	Env     []string
	Prepare func() error
	Os      map[string]GoPortal
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
	if !g.prepare() {
		return
	}
	g.run()
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
	if goos, ok := g.Os[exec.GetEnv("GOOS")]; ok {
		g.Env = append(g.Env, goos.Env...)
		g.LdFlags = append(g.LdFlags, goos.LdFlags...)
	}

	target := fmt.Sprintf("./cmd/%s", g.Target)
	g.Args = append([]string{g.Cmd}, g.Args...)
	g.Args = g.arg("-tags", g.Tags)
	g.Args = g.arg("-ldflags", strings.Join(g.LdFlags, " "))
	g.Args = append(g.Args, target)
	log.Printf("$ %s", strings.Join(g.Args, " "))
	cmd := exec.Cmd{Cmd: "go", Args: g.Args}.Default().AddEnv().Build()
	if err := cmd.Run(); err != nil {
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
