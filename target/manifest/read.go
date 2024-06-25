package manifest

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

func Read(src fs.FS) (p target.Manifest, err error) {
	err = Load(&p, src, target.PortalJsonFilename)
	return
}

func Load(manifest *target.Manifest, src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, manifest)
	return
}
