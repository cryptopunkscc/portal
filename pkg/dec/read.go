package dec

import (
	"io/fs"
)

type Unmarshal func(in []byte, out interface{}) (err error)

type Unmarshaler interface {
	Load(dst any, src fs.FS, name string) (err error)
}

type Load func(dst any, src fs.FS, name string) (err error)

func (l Load) Load(dst any, src fs.FS, name string) (err error) { return l(dst, src, name) }

type Unmarshalers map[string]Unmarshal

var _ Unmarshaler = Unmarshalers{}

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
	if err == nil {
		err = fs.ErrNotExist
	}
	return err
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
