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
	Portal
	IsProject()
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

type Html interface {
	IsHtml()
}

type AppHtml interface {
	App
	Html
}

type PortalHtml interface {
	Portal
	Html
}

type ProjectHtml interface {
	ProjectNodeModule
	Html
	DistHtml() DistHtml
}

type DistHtml interface {
	Dist
	Html
}

type BundleHtml interface {
	Bundle
	Html
}

type Js interface {
	IsJs()
}

type ProjectJs interface {
	ProjectNodeModule
	Js
	DistJs() DistJs
}

type AppJs interface {
	App
	Js
}

type PortalJs interface {
	Portal
	Js
}

type DistJs interface {
	Dist
	Js
}

type BundleJs interface {
	Bundle
	Js
}

type Exec interface {
	Executable() Source
}

type AppExec interface {
	App
	Exec
}

type DistExec interface {
	Dist
	Exec
}

type BundleExec interface {
	Bundle
	Exec
}

type ProjectGo interface {
	Project
	DistGolang() DistExec
}
