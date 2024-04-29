package runner

import (
	"io/fs"
)

type Target interface {
	Path() string
	Files() fs.FS
	Type() Type
}

type Type int

const (
	Invalid Type = iota
	Backend
	Frontend
)
