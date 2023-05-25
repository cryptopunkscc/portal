package main

import (
	"io"
	"os"
	"path/filepath"
)

var _ FileStore = &DirStore{}

type DirStore struct {
	root string
}

func NewDirStore(root string) *DirStore {
	return &DirStore{root: root}
}

func (store *DirStore) Open(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(store.root, path))
}
