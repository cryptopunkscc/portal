package appstore

import (
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"path/filepath"
)

func Install(src string) (err error) {
	file, err := source.File(src)
	if err != nil {
		return err
	}
	for _, t := range target.List(
		apps.Resolver[target.Bundle_](),
		file,
	) {
		src = t.Abs()
		dst := filepath.Join(portalAppsDir, filepath.Base(t.Abs()))

		log.Printf("Installing %s to %s", src, dst)
		err = fs2.CopyFile(src, dst)
		if err != nil {
			log.Printf("Error copying file %s: %v", src, err)
			return
		}
	}
	return
}
