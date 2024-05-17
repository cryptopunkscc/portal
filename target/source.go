package target

import (
	"io/fs"
)

type Type int

func (t Type) Is(p Type) bool {
	return t&p == p
}

const (
	TypeNone     = Type(0x0)
	TypeBackend  = Type(0x1)
	TypeFrontend = Type(0x2)
	TypeDev      = Type(0x4)
	TypeBundle   = Type(0x8)
	TypeAll      = TypeBackend | TypeFrontend | TypeDev | TypeBundle
)

type Source interface {
	Abs() string
	Path() string
	Files() fs.FS
	Type() Type
	Lift() Source
	Parent() Source
	IsFile() bool
}

type Portals[T Portal] map[string]T

type Portal interface {
	Source
	Manifest() *Manifest
}

type Template interface {
	Source
	Info() TemplateInfo
	Name() string
}

type NodeModule interface {
	Source
	PkgJson() *PackageJson
}

type Project interface {
	IsProject()
	NodeModule
	Portal
}

type App interface {
	IsApp()
	Portal
}

type Bundle interface {
	IsBundle()
	App
}

type Dist interface {
	IsDist()
	App
}

type Frontend interface {
	IsFrontend()
}

type AppFrontend interface {
	App
	Frontend
}

type PortalFrontend interface {
	Portal
	Frontend
}

type ProjectFrontend interface {
	Project
	Frontend
}

type DistFrontend interface {
	Dist
	Frontend
}

type BundleFrontend interface {
	Bundle
	Frontend
}

type Backend interface {
	IsBackend()
}

type ProjectBackend interface {
	Project
	Backend
}

type AppBackend interface {
	App
	Backend
}

type PortalBackend interface {
	Portal
	Backend
}

type DistBackend interface {
	Dist
	Backend
}

type BundleBackend interface {
	Bundle
	Backend
}
