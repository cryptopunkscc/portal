//go:build mage

package main

import (
	"github.com/cryptopunkscc/portal/api"
	"github.com/cryptopunkscc/portal/runner/apps_build"
	"github.com/cryptopunkscc/portal/runner/astrald_build"
	"github.com/cryptopunkscc/portal/runner/js_build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"runtime"
	"strings"
)

var Aliases = map[string]interface{}{
	"apps": Install.Apps,
	"a":    Install.Astrald,
	"p":    Install.Portald,
	"c":    Install.Cli,
	"all":  Install.All,
	"ir":   Build.Installer,
}

var goos = []string{runtime.GOOS}

func Goos(str string) {
	if len(goos) > 0 {
		goos = strings.Split(str, " ")
	}
}

type Install mg.Namespace

func (Install) Astrald() error {
	return sh.RunV("go", "install", "github.com/cryptopunkscc/astrald/cmd/astrald@"+api.AstralVersion)
}

func (Install) Portald() error {
	return sh.RunV("go", "install", "./cmd/portald")
}

func (Install) Cli() error {
	return sh.RunV("go", "install", "./cmd/portal")
}

func (Install) Apps() (err error) {
	panic("fixme")
	//results, err := install.Runner{
	//	AppsDir: env.PortaldApps.Get(),
	//	Tokens: token.Repository{
	//		Dir: env.PortaldTokens.Get(),
	//	},
	//}.BundlesBySource("astrald")
	//if err != nil {
	//	return
	//}
	//for result := range results {
	//	println(fmt.Println(result))
	//}
	//return
}

func (Install) All() {
	mg.Deps(
		Install.Astrald,
		Install.Portald,
		Install.Cli,
		Install.Apps,
	)
}

type Build mg.Namespace

func (Build) Installer() error {
	mg.Deps(
		Build.Astrald,
		Build.Portald,
		Build.Cli,
		Build.Apps,
	)
	return sh.RunV("go", "build", "-o", "./bin", "./cmd/portal-installer/")
}

func (Build) Apps() (err error) {
	mg.Deps(
		Build.JsLib,
	)
	return apps_build.Run()
}

func (Build) Cli() error {
	return sh.RunV("go", "build", "-o", "./cmd/portal-installer/bin/", "./cmd/portal")
}

func (Build) Astrald() error {
	return astrald_build.Run()
}

func (Build) Portald() error {
	return sh.RunV("go", "build", "-o", "./cmd/portal-installer/bin/", "./cmd/portald")
}

func (Build) JsLib() error {
	if changed, err := target.Path(
		"./core/js/embed/",
		"./core/js/src/",
		"./core/js/all.js",
		"./core/js/common.js",
		"./core/js/package.json",
	); err != nil {
		return err
	} else if !changed {
		return nil
	}
	return js_build.Run()
}
