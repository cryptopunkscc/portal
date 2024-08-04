package dec

import (
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

type Unmarshal func(in []byte, out interface{}) (err error)

type Unmarshalers map[string]Unmarshal

func From(unmarshalers ...Unmarshalers) (out Unmarshalers) {
	out = make(Unmarshalers)
	for _, unmarshaler := range unmarshalers {
		for ext, unmarshal := range unmarshaler {
			out[ext] = unmarshal
		}
	}
	return
}

func (u Unmarshalers) Load(dst any, src fs.FS, name string) (err error) {
	for ext, unmarshal := range u {
		if err = load(unmarshal, dst, src, name, ext); err == nil {
			return
		}
	}
	return target.ErrNotTarget
}

func load(
	unmarshal Unmarshal,
	dst any,
	src fs.FS,
	name string,
	ext string,
) (err error) {
	name = name + "." + ext
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	return unmarshal(file, dst)
}
