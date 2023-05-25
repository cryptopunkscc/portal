package main

import (
	"astraljs"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/data"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var _ http.Handler = &App{}

//go:embed apphost_wails.js
var apphostWails string

type Manifest struct {
	Name string `json:"name"`
}

type App struct {
	Title string
	FileStore
}

func NewApp(appName string) (*App, error) {
	if _, err := data.Parse(appName); err == nil {
		panic("astral storage not implemented")
	}

	var app = &App{}
	var path = appName

	if path[0] != '/' {
		var cwd, _ = os.Getwd()
		path = filepath.Join(cwd, appName)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	app.Title = filepath.Base(path)

	switch {
	case info.IsDir():
		app.FileStore = NewDirStore(path)

	case info.Mode().IsRegular():
		filetype, err := fileType(path)
		if err != nil {
			return nil, err
		}
		switch filetype {
		case "text/html":
			var mem = NewMemStore()
			mem.Entries["index.html"], _ = os.ReadFile(path)
			mem.Entries["apphost.js"] = []byte(apphostWails + astraljs.AppHostJsClient())
			app.FileStore = mem

		case "application/zip":
			app.FileStore, err = NewZipStore(path)
			if err != nil {
				return nil, err
			}

		default:
			return nil, errors.New("unsupported file type: " + filetype)
		}

	default:
		return nil, errors.New("app must be a file or a directory")
	}

	if f, err := app.Open("manifest.json"); err == nil {
		var jdec = json.NewDecoder(f)
		var m Manifest
		if err := jdec.Decode(&m); err == nil {
			app.Title = m.Name
		}
	}

	return app, nil
}

func fileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := strings.Split(http.DetectContentType(buffer), ";")[0]

	return contentType, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var filename = strings.TrimPrefix(req.URL.Path, "/")

	if filename == "" {
		filename = "index.html"
	}

	if r, err := app.Open(filename); err != nil {
		fmt.Fprintf(os.Stderr, "GET /%s 404\n", filename)
		w.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprintf(os.Stdout, "GET /%s 200\n", filename)
		w.WriteHeader(http.StatusOK)
		io.Copy(w, r)
	}
}
