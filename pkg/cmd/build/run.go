package build

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
	"os"
	"path"
	"sync"
)

func Run(dir string) (err error) {
	r, err := runner.New(dir, runner.DevTargets)
	if err != nil {
		return
	}
	return RunRunner(r)
}

func RunRunner(r *runner.Runner) (err error) {
	targets := append(r.Backends, r.Frontends...)

	wait := sync.WaitGroup{}
	wait.Add(len(targets))
	for _, target := range targets {
		go func(src string) {
			defer wait.Done()
			if err := npmInstall(src); err != nil {
				log.Println(err)
			}
			if err := npmRunBuild(src); err != nil {
				log.Println(err)
			}
			if err := copyManifest(src); err != nil {
				log.Println(err)
			}
		}(target.Path)
	}
	wait.Wait()
	return
}

func npmInstall(dir string) error {
	return exec.Run(dir, "npm", "install")
}

func npmRunBuild(dir string) error {
	return exec.Run(dir, "npm", "run", "build")
}

func copyManifest(src string) (err error) {
	manifest := bundle.Base(src)
	_ = manifest.LoadPath(src, "package.json")
	_ = manifest.LoadPath(src, bundle.PortalJson)
	if manifest.Icon != "" {
		iconSrc := path.Join(src, manifest.Icon)
		iconName := "icon" + path.Ext(manifest.Icon)
		iconDst := path.Join(src, "dist", iconName)
		if err = fs.CopyFile(iconSrc, iconDst); err != nil {
			return
		}
		manifest.Icon = iconName
	}

	bytes, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(src, "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
