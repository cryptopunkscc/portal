package main

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	"strings"
)

var GoWails = GoPortal{
	Tags:    "desktop,wv2runtime.download,production,webkit2_41",
	LdFlags: "-w -s",
}

var GoWailsDev = GoPortal{
	Tags: "dev,webkit2_41",
}

func goPortal() GoPortal         { return GoPortal{}.target("portal") }
func goPortalApp() GoPortal      { return GoPortal{}.target("portal-app") }
func goPortalAppWails() GoPortal { return GoWails.target("portal-app-wails") }
func goPortalAppGoja() GoPortal  { return GoPortal{}.target("portal-app-goja") }
func goPortalTray() GoPortal     { return GoPortal{}.target("portal-tray") }
func goPortalDev() GoPortal      { return GoPortal{}.target("portal-dev") }
func goPortalDevExec() GoPortal  { return GoPortal{}.target("portal-dev-exec") }
func goPortalDevWails() GoPortal { return GoWailsDev.target("portal-dev-wails") }
func goPortalDevGoja() GoPortal  { return GoPortal{}.target("portal-dev-goja") }
func goPortalDevGo() GoPortal    { return GoPortal{}.target("portal-dev-go") }

func buildPortalInstaller() { GoPortal{Out: "./bin/"}.target("portal-installer").Build() }

type GoPortal struct {
	Cmd     string
	Tags    string
	LdFlags string
	Out     string
	Target  string
	Args    []string
}

func (g GoPortal) Install() {
	g.Cmd = "install"
	g.run()
}

func (g GoPortal) Build() {
	g.Cmd = "build"
	if g.Out == "" {
		g.Out = "./cmd/portal-installer/bin/"
	}
	g.Args = g.arg("-o", g.Out)
	g.run()
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
