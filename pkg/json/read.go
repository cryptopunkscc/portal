package json

import (
	"encoding/json"
	"io/fs"
)

func Read[T any](src fs.FS, name string) (p T, err error) {
	err = Load(&p, src, name)
	return
}

func Load(dst any, src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, dst)
	return
}
