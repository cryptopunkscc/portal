package main

import (
	"bytes"
	"errors"
	"io"
)

var _ FileStore = &MemStore{}

type MemStore struct {
	Entries map[string][]byte
}

type MemReadCloser struct {
	io.Reader
}

func NewMemStore() *MemStore {
	return &MemStore{
		Entries: map[string][]byte{},
	}
}

func (MemReadCloser) Close() error { return nil }

func (store *MemStore) Open(path string) (io.ReadCloser, error) {
	if data, found := store.Entries[path]; found {
		return MemReadCloser{bytes.NewReader(data)}, nil
	}

	return nil, errors.New("not found")
}
