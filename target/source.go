package target

import (
	"io/fs"
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
	Portal
	Dist() Dist
}

type ProjectNodeModule interface {
	Project
	NodeModule
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
	ProjectNodeModule
	Frontend
	DistFrontend() DistFrontend
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
	ProjectNodeModule
	Backend
	DistBackend() DistBackend
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
