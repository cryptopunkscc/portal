//go:build mage

package main

import (
	"errors"
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/js"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/gpg"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	plog.Verbosity = 100
}

var Aliases = map[string]interface{}{
	"apps": Install.Apps,
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

var clean = false

func Clean() { clean = true }

type Install mg.Namespace

func (Install) Astrald() (err error) {
	d, err := golang.ProjectDependency("astrald")
	if err != nil {
		return
	}
	return d.Install("cmd/astrald")
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

var out string

func (Build) Out(dir string) { out = dir }

func (Build) Installer() error {
	defer clearVersion()
	defer gpgSignPortalInstallers()
	cleanInstallerBin()
	resolveVersion()
	mg.Deps(
		Build.Astrald,
		Build.Portald,
		Build.Cli,
		Build.Apps,
	)
	o := "./bin/"
	if len(out) > 0 {
		o = out
	}
	return sh.RunV("go", "build", "-o", o, "./cmd/install-portal-to-astral/")
}

func (Build) Apps() (err error) {
	mg.Deps(
		Build.JsLib,
	)
	var args []string
	args = append(args, "pack")
	if clean {
		args = append(args, "clean")
	}
	return apps.Build(args...)
}

func (Build) Cli() error {
	return sh.RunV("go", "build", "-o", "./cmd/install-portal-to-astral/bin/", "./cmd/portal")
}

func (Build) Astrald() error {
	n := "cryptopunkscc/astrald"
	c := "./cmd/astrald"
	o := "cmd/install-portal-to-astral/bin/"
	d, err := golang.ProjectDependency(n)
	if err != nil {
		return err
	}
	if d.Replace == "" {
		if err = d.Get(); err != nil {
			return err
		}
	}
	o = filepath.Join(d.Dir, o)
	return d.Build(c, o)
}

func (Build) Portald() error {
	return sh.RunV("go", "build", "-o", "./cmd/install-portal-to-astral/bin/", "./cmd/portald")
}

func (Build) JsLib() error {
	if !clean {
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
	}
	return js.BuildPortalLib()
}

func resolveVersion() {
	file, err := os.Create("./api/version/name")
	if err != nil {
		return
	}
	defer file.Close()
	name := version.Resolve()
	if _, err = file.WriteString(name); err != nil {
		panic(err)
	}
}

func clearVersion() {
	file, err := os.Create("./api/version/name")
	if err != nil {
		panic(err)
	}
	if err = file.Truncate(0); err != nil {
		panic(err)
	}
}

func cleanInstallerBin() {
	path := "./cmd/install-portal-to-astral/bin"
	_ = os.RemoveAll(path)
	err := os.Mkdir(path, 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}
}

func gpgSignPortalInstallers() {
	_ = filepath.WalkDir("./bin", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() ||
			!strings.HasPrefix(d.Name(), "install-portal-to-astral") ||
			strings.HasSuffix(d.Name(), ".sig") {
			return nil
		}
		time.Sleep(1 * time.Second)
		gpg.Sign(path)
		return nil
	})
}
