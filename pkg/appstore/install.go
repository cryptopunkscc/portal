package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"path"
)

func Install(src string) (err error) {
	for t := range project.FindInPath[target.Bundle](src) {
		src = t.Abs()
		dst := path.Join(portalAppsDir, path.Base(t.Path()))

		err = fs.CopyFile(src, dst)
		log.Printf("Installing %s to %s", src, dst)
		if err != nil {
			log.Printf("Error copying file %s: %v", src, err)
			return
		}
	}
	return
}
