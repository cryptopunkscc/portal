package bundle

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
)

type Manifest struct {
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Package     string `json:"package,omitempty"`
	Version     string `json:"version,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

func Base(src string) Manifest {
	return Manifest{
		Name:    path.Base(src),
		Version: "0.0.0",
	}
}

func ReadManifestFs(src fs.FS) (p *Manifest, err error) {
	pp := Manifest{}
	if err = pp.LoadFs(src, "portal.json"); err == nil {
		p = &pp
	}
	return
}

func (m *Manifest) LoadFs(src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	log.Println("ReadManifest: ", string(file))
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &m)
	log.Printf("ReadManifest.struct: %v", m)
	return
}

func (m *Manifest) LoadPath(src string, name string) (err error) {
	bytes, err := os.ReadFile(path.Join(src, name))
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, &m)
}
