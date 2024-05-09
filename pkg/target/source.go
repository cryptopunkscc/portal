package target

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"io/fs"
)

type Type int

func (t Type) Has(p Type) bool {
	return t&p == p
}

const (
	Invalid Type = iota
	Backend
	Frontend
)

type Source interface {
	Path() string
	Files() fs.FS
	Type() Type
}

type Portals[T Portal] map[string]T

type Portal interface {
	Source
	Manifest() bundle.Manifest
}

type NodeModule interface {
	Source
	PkgJson() bundle.PackageJson
}

type Project interface {
	NodeModule
	Portal
}

type App interface {
	Source
	Portal
	App() // required to disable Project to App type casting
}
