package build

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
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
