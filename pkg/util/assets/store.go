package assets

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Store interface {
	Open(string) (io.ReadCloser, error)
}

type OverlayStore struct {
	Stores []Store
}

type FsStore struct {
	fs.FS
}

func (fss *FsStore) Open(path string) (rc io.ReadCloser, err error) {
	rc, err = fss.FS.Open(path)
	return
}

func (o *OverlayStore) Open(s string) (io.ReadCloser, error) {
	for _, store := range o.Stores {
		if rc, err := store.Open(s); err == nil {
			return rc, nil
		}
	}
	return nil, errors.New("not found")
}

func SingleFileStore(path string, name string) (s Store, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return
	}
	return &MemStore{map[string][]byte{name: file}}, nil
}

type MemStore struct {
	Entries map[string][]byte
}

type MemReadCloser struct {
	io.Reader
}

func (MemReadCloser) Close() error { return nil }

func (store *MemStore) Open(path string) (io.ReadCloser, error) {
	if data, found := store.Entries[path]; found {
		return MemReadCloser{bytes.NewReader(data)}, nil
	}

	return nil, errors.New("not found")
}

type ZipStore struct {
	zip *zip.ReadCloser
}

func NewZipStore(zipPath string) (*ZipStore, error) {
	var err error
	var store = &ZipStore{}

	store.zip, err = zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *ZipStore) Open(s string) (io.ReadCloser, error) {
	return store.zip.Open(s)
}

type DirStore struct {
	root string
}

func NewDirStore(root string) *DirStore {
	return &DirStore{root: root}
}

func (store *DirStore) Open(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(store.root, path))
}
