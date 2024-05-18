package target

import (
	"io/fs"
)

type Type int

func (t Type) Is(p Type) bool {
	return t&p == p
}

const (
	None         = Type(0x0)
	TypeBackend  = Type(0x1)
	TypeFrontend = Type(0x2)
	TypeDev      = Type(0x4)
	TypeBundle   = Type(0x8)
)

type Source interface {
	Abs() string
	Path() string
	Files() fs.FS
	Type() Type
	Lift() Source
	Parent() Source
}

type Portals[T Portal] map[string]T

type Portal interface {
	Source
	Manifest() *Manifest
}

type NodeModule interface {
	Source
	PkgJson() *PackageJson

	IsPortalLib() bool
	CanNpmRunBuild() bool
}

type Project interface {
	Project()
	NodeModule
	Portal
}

type App interface {
	App()
	Source
	Portal
}

type Bundle interface {
	Bundle()
	App
}

type Dist interface {
	Dist()
	App
}
