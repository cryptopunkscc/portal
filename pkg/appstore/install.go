package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
	"path"
)

func Install(src string) (err error) {
	for target := range project.Find[project.Bundle](os.DirFS(src), ".") {
		wd := ""
		wd, err = os.Getwd()
		if err != nil {
			return
		}
		src = path.Join(wd, target.Path())
		dst := path.Join(portalAppsDir, path.Base(target.Path()))

		err = fs.CopyFile(src, dst)
		log.Printf("Installing %s to %s", src, dst)
		if err != nil {
			log.Printf("Error copying file %s: %v", src, err)
			return
		}
	}
	return
}
