package main

import (
	"errors"
	"io"
)

var _ FileStore = &OverlayStore{}

type OverlayStore struct {
	Stores []FileStore
}

func NewOverlayStore(stores ...FileStore) *OverlayStore {
	return &OverlayStore{Stores: stores}
}

func (o *OverlayStore) Open(s string) (io.ReadCloser, error) {
	for _, store := range o.Stores {
		if rc, err := store.Open(s); err == nil {
			return rc, nil
		}
	}
	return nil, errors.New("not found")
}
