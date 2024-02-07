package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
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
	FileStore
	Title    string
	tempFile string
}

func NewApp(appName string) (*App, error) {
	var app = &App{}

	//if blockID, err := data.Parse(appName); err == nil {
	//	//TODO: don't download the block, read it remotely
	//	file, err := downloadBlock(blockID)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	appName = file
	//	app.tempFile = file
	//}

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
		app.FileStore = NewOverlayStore(NewDirStore(path), SDKStore)

	case info.Mode().IsRegular():
		filetype, err := fileType(path)
		if err != nil {
			return nil, err
		}
		switch filetype {
		case "text/html":
			var store = NewMemStore()
			store.Entries["index.html"], _ = os.ReadFile(path)
			app.FileStore = NewOverlayStore(store, SDKStore)

		case "application/zip":
			store, err := NewZipStore(path)
			if err != nil {
				return nil, err
			}
			app.FileStore = NewOverlayStore(store, SDKStore)

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

//func downloadBlock(blockID data.ID) (string, error) {
//	storage := astral.Client.LocalStorage()
//	r, err := storage.Read(blockID, 0, 0)
//	if err != nil {
//		return "", err
//	}
//
//	file, err := os.CreateTemp("", ".astral.app-*")
//	if err != nil {
//		return "", err
//	}
//
//	_, err = io.Copy(file, r)
//	if err != nil {
//		return "", err
//	}
//
//	return file.Name(), nil
//}

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

func (app *App) Cleanup() {
	if app.tempFile != "" {
		os.Remove(app.tempFile)
	}
}
