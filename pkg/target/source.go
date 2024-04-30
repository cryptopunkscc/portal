package target

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"io/fs"
)

type Source interface {
	Path() string
	Files() fs.FS
	Type() Type
}

type Type int

func (t Type) Has(p Type) bool {
	return t&p == p
}

const (
	Invalid Type = iota
	Backend
	Frontend
)

type Project interface {
	NodeModule
	App
}

type NodeModule interface {
	Source
	PkgJson() bundle.PackageJson
}

type App interface {
	Source
	Manifest() bundle.Manifest
}
