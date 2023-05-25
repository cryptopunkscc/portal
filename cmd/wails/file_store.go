package main

import (
	"io"
)

type FileStore interface {
	Open(string) (io.ReadCloser, error)
}
