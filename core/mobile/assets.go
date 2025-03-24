package core

import (
	"github.com/cryptopunkscc/portal/api/mobile"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type assets struct{ files fs.FS }

func (a assets) Get(uri string) (out mobile.Asset, err error) {
	uri = filepath.Clean(uri)

	mimetype := mime.TypeByExtension(filepath.Ext(uri))
	if len(mimetype) == 0 {
		if mimetype, err = a.mimeFromFile(uri); err == nil {
			mimetype = "application/octet-stream; charset=utf-8"
		}
	}

	encoding := "UTF-8"
	chunks := strings.Split(mimetype, ";")
	mimetype = chunks[0]
	if len(chunks) > 1 {
		encoding = chunks[1]
		encoding = strings.TrimPrefix(encoding, " charset=")
		encoding = strings.ToUpper(encoding)
	}

	file, err := a.files.Open(uri)
	if err != nil {
		return
	}
	out = asset{
		//mimetype: strings.Split(mime.String(), ";")[0],
		encoding: encoding,
		mimetype: mimetype,
		Reader:   reader{file},
		Closer:   file,
	}
	return
}

func (a assets) mimeFromFile(uri string) (mime string, err error) {
	file, err := a.files.Open(uri)
	if err != nil {
		return
	}
	contentTypeBuff := make([]byte, 512)
	n, err := file.Read(contentTypeBuff)
	if err != nil {
		return
	}
	_ = file.Close()
	if n < 512 {
		contentTypeBuff = contentTypeBuff[:n]
	}
	mime = http.DetectContentType(contentTypeBuff)
	return
}

type asset struct {
	mimetype string
	encoding string
	mobile.Reader
	mobile.Closer
}

func (a asset) Data() mobile.ReadCloser { return a }
func (a asset) Encoding() string        { return a.encoding }
func (a asset) Mime() string            { return a.mimetype }
