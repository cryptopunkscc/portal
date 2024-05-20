package portal

import (
	"encoding/json"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

func ReadManifest(src fs.FS) (p target.Manifest, err error) {
	err = LoadManifest(&p, src, target.PortalJsonFilename)
	return
}

func LoadManifest(m *target.Manifest, src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, m)
	return
}
