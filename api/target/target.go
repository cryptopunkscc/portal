package target

import (
	"io/fs"
)

type Source interface {
	Abs() string
	Path() string
	File() (fs.File, error)
	FS() fs.FS
	Sub(dir string) (Source, error)
	IsDir() bool
}

type Template interface {
	Source
	Info() TemplateInfo
	Name() string
}

type Portal_ interface {
	Source
	Manifest() *Manifest
	MarshalJSON() ([]byte, error)
}

type Portal[T any] interface {
	Portal_
	Target() T
}

type Portals[T Portal_] []T

type NodeModule interface {
	Source
	PkgJson() *PackageJson
}

type Project_ interface {
	Portal_
	Build() Builds
	Dist_() Dist_
}

type Project[T any] interface {
	Portal[T]
	Build() Builds
	Dist_() Dist_
	Dist() Dist[T]
}

type ProjectNpm_ interface {
	Project_
	PkgJson() *PackageJson
}

type ProjectNpm[T any] interface {
	Project[T]
	PkgJson() *PackageJson
}

type App_ interface {
	Portal_
	IsApp()
}

type App[T any] interface {
	Portal[T]
	IsApp()
}

type Bundle interface {
	Source
	Package() Source
}

type Bundle_ interface {
	Package() Source
	App_
}

type AppBundle[T any] interface {
	Package() Source
	App[T]
}

type Dist_ interface {
	IsDist()
	App_
}

type Dist[T any] interface {
	IsDist()
	App[T]
}

type Html interface{ IndexHtml() }
type PortalHtml Portal[Html]
type AppHtml App[Html]
type DistHtml Dist[Html]
type BundleHtml AppBundle[Html]
type ProjectHtml ProjectNpm[Html]

type Js interface{ MainJs() }
type PortalJs Portal[Js]
type AppJs App[Js]
type DistJs Dist[Js]
type BundleJs AppBundle[Js]
type ProjectJs ProjectNpm[Js]

type Exec interface{ Executable() Source }
type PortalExec Portal[Exec]
type AppExec App[Exec]
type DistExec Dist[Exec]
type BundleExec AppBundle[Exec]
type ProjectExec Project[Exec]

type ProjectGo interface {
	Project[Exec]
	IsGo()
}
