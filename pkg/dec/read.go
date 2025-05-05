package dec

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io/fs"
)

type Unmarshal func(in []byte, out any) (err error)

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

func (u Unmarshalers) Unmarshal(in []byte, out any) (err error) {
	defer plog.TraceErr(&err)
	for _, unmarshal := range u {
		if err = unmarshal(in, out); err == nil {
			return
		}
	}
	return fs.ErrNotExist
}

func (u Unmarshalers) Load(dst any, src fs.FS, names ...string) (err error) {
	defer plog.TraceErr(&err)
	for _, name := range names {
		for ext, unmarshal := range u {
			if err = load(unmarshal, dst, src, name, ext); err == nil {
				return
			}
		}
	}
	return fs.ErrNotExist
}

func load(
	unmarshal Unmarshal,
	dst any,
	src fs.FS,
	name string,
	ext string,
) (err error) {
	name = name + "." + ext
	b, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	return unmarshal(b, dst)
}
