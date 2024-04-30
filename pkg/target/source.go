package target

import (
	"io/fs"
)

type Source interface {
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
