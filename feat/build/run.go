package build

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
	"path"
	"sync"
)

func Run(dir string) (err error) {
	targets := project.DevTargets(os.DirFS(dir))
	wait := sync.WaitGroup{}
	for target := range targets {
		wait.Add(1)
		go func(target project.PortalNodeModule) {
			defer wait.Done()
			if err := target.NpmInstall(); err != nil {
				log.Println(err)
				return
			}
			if err := target.NpmRunBuild(); err != nil {
				log.Println(err)
				return
			}
			if err := target.CopyManifest(); err != nil {
				log.Println(err)
				return
			}
		}(target)
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
