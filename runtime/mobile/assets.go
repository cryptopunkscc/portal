package runtime

import (
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/wailsapp/mimetype"
	"io/fs"
)

type assets struct{ files fs.FS }

func (a assets) Get(uri string) (out mobile.Asset, err error) {
	file, err := a.files.Open(uri)
	if err != nil {
		return
	}
	out = asset{reader{file}, file}
	return
}

type asset struct {
	mobile.Reader
	mobile.Closer
}

func (a asset) Data() mobile.ReadCloser { return a }
func (a asset) Encoding() string        { return "UTF-8" }
func (a asset) Mime() string {
	if r, err := mimetype.DetectReader(a); err == nil {
		return r.String()
	}
	return ""
}
