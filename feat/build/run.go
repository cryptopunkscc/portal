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
	targets, err := runner.DevTargets(dir)
	if err != nil {
		return err
	}
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
			if err := CopyManifest(src); err != nil {
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

func CopyManifest(src string) (err error) {
	m := bundle.Base(src)
	_ = m.LoadPath(src, "package.json")
	_ = m.LoadPath(src, bundle.PortalJson)
	if m.Icon != "" {
		iconSrc := path.Join(src, m.Icon)
		iconName := "icon" + path.Ext(m.Icon)
		iconDst := path.Join(src, "dist", iconName)
		if err = fs.CopyFile(iconSrc, iconDst); err != nil {
			return
		}
		m.Icon = iconName
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(src, "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
