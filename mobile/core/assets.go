package core

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/bind/js/embed/android"
)

type assets struct{ files fs.FS }

func (a assets) Get(uri string) (out mobile.Asset, err error) {
	uri = filepath.Clean(uri)

	if strings.HasSuffix(uri, "portal.js") {
		out = asset{
			encoding: "UTF-8",
			mimetype: "application/javascript",
			Reader:   reader{strings.NewReader(android.JsString)},
			Closer:   io.NopCloser(nil),
		}
		return
	}

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

type reader struct{ io.Reader }

func (r reader) Read(arr []byte) (n int, err error) {
	n, err = r.Reader.Read(arr)
	return
}

func (r reader) ReadN(n int) (arr []byte, err error) {
	var l int
	arr = make([]byte, n)
	if l, err = r.Reader.Read(arr); err == nil {
		arr = arr[:l]
	}
	return
}

func (r reader) ReadAll() (all []byte, err error) {
	all, err = io.ReadAll(r)
	return all, err
}
