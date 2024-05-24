package manifest

import (
	"encoding/json"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
)

func Read(src fs.FS) (p target.Manifest, err error) {
	err = Load(&p, src, target.PortalJsonFilename)
	return
}

func Load(m *target.Manifest, src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, m)
	return
}