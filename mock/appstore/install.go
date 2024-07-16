package appstore

import (
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"log"
	"path"
)

func Install(src string) (err error) {
	for _, t := range apps.FromPath[target.Bundle](src) {
		src = t.Abs()
		dst := path.Join(portalAppsDir, path.Base(t.Path()))

		err = fs2.CopyFile(src, dst)
		log.Printf("Installing %s to %s", src, dst)
		if err != nil {
			log.Printf("Error copying file %s: %v", src, err)
			return
		}
	}
	return
}
