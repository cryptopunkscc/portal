package main

import (
	"errors"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/gpg"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (d *Install) buildInstaller(platforms ...string) {

	if len(platforms) == 0 {
		platforms = []string{runtime.GOOS}
	}

	for _, platform := range platforms {
		cleanInstallerBin()
		goos := "GOOS=" + platform
		log.Println(goos)
		exec.SetEnv(goos)
		d.buildEmbedApps(platform)
		d.buildInstallerFor()
		exec.SetEnv()
	}

	gpgSignPortalInstallers()
}

func (d *Install) buildInstallerFor() {
	mkdirBin()
	buildAstral()
	buildAnc()
	goPortal.Build()
	goPortalApp.Build()
	goPortalApps.Build()
	goPortalAppGoja.Build()
	goPortalAppWails.Build()
	goPortalCreate.Build()
	goPortalList.Build()
	goPortalBuild.Build()
	goPortalDev.Build()
	goPortalDevGo.Build()
	goPortalDevGoja.Build()
	goPortalDevExec.Build()
	goPortalDevWails.Build()
	buildPortalInstaller()
}

func cleanInstallerBin() {
	_ = os.RemoveAll("./cmd/portal-installer/bin")
}

func mkdirBin() {
	err := os.Mkdir("./cmd/portal-installer/bin/", 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		panic(err)
	}
}

func gpgSignPortalInstallers() {
	_ = filepath.WalkDir("./bin", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() ||
			!strings.HasPrefix(d.Name(), "portal-installer") ||
			strings.HasSuffix(d.Name(), ".sig") {
			return nil
		}
		time.Sleep(1 * time.Second)
		gpg.Sign(path)
		return nil
	})
}
