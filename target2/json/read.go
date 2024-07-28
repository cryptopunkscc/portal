package json

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/target2"
	"io/fs"
)

func Read[T any](src target2.Source, name string) (p T, err error) {
	err = Load(&p, src, name)
	return
}

func Load(dst any, src target2.Source, name string) (err error) {
	file, err := fs.ReadFile(src.Files(), name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, dst)
	return
}
