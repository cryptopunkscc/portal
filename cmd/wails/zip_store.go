package main

import (
	"archive/zip"
	"io"
)

var _ FileStore = &ZipStore{}

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
