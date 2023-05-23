package wails

import (
	"astraljs"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//go:embed apphost_wails.js
var apphostWails string

type FileLoader struct {
	app astraljs.WebApp
}

func NewFileLoader(app astraljs.WebApp) *FileLoader {
	return &FileLoader{
		app: app,
	}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	path := strings.TrimPrefix(req.URL.Path, "/")
	switch path {
	case "", "index.html":
		data = []byte(h.app.Source)
	case "apphost.js":
		data = []byte(apphostWails + astraljs.AppHostJsClient())
	default:
		err = fmt.Errorf("[astral-runtime-wails] unhandled file: %s", path)
	}

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte(fmt.Sprintf("[astral-runtime-wails] could not load file %s", path)))
		log.Fatalln(err)
		return
	}

	_, _ = res.Write(data)
}
