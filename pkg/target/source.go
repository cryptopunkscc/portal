package target

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"io/fs"
)

type Type int

func (t Type) Is(p Type) bool {
	return t&p == p
}

const (
	None     = Type(0x0)
	Backend  = Type(0x1)
	Frontend = Type(0x2)
	Dev      = Type(0x4)
	Bundle   = Type(0x8)
)

type Source interface {
	Path() string
	Abs() string
	Files() fs.FS
	Type() Type
	Parent() Source
}

type Portals[T Portal] map[string]T

type Portal interface {
	Source
	Manifest() bundle.Manifest
}

type NodeModule interface {
	Source
	PkgJson() bundle.PackageJson

	// TODO temprary
	CanNpmRunBuild() bool
	NpmInstall() (err error)
	NpmRunBuild() (err error)
	InjectDependencies(modules []NodeModule) (err error)
	IsPortalLib() bool
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
