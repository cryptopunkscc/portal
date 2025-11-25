package target

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/api/manifest"
)

type Source interface {
	Abs() string
	Path() string
	File() (fs.File, error)
	FS() fs.FS
	Sub(path ...string) (Source, error)
	CopyTo(path ...string) (err error)
	IsDir() bool
}

type Template interface {
	Source
	Info() TemplateInfo
	Name() string
}

type Portal_ interface {
	Source
	Api() *manifest.Api
	Config() *manifest.Config
	Manifest() *manifest.App
	MarshalJSON() ([]byte, error)
}

type Portal[T any] interface {
	Portal_
	Runtime() T
}

type Portals[T Portal_] []T

type NodeModule interface {
	Source
	PkgJson() *PackageJson
}

type Project_ interface {
	Portal_
	Build() *manifest.Builds
	Dist_(platform ...string) Dist_
	Changed() bool
}

type Project[T any] interface {
	Project_
	Portal[T]
	Dist(platform ...string) Dist[T]
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
	Release() *manifest.Release
	Version() string
}

type App[T any] interface {
	App_
	IsApp()
	Runtime() T
}

type Bundle interface {
	Source
	Package() Source
}

type Dist_ interface {
	IsDist()
	App_
}

type Bundle_ interface {
	Package() Source
	Dist_
}

type Dist[T any] interface {
	IsDist()
	App[T]
}

type AppBundle[T any] interface {
	Package() Source
	Dist[T]
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

type Exec interface {
	Target() manifest.Target
	Executable() Source
}
type PortalExec Portal[Exec]
type AppExec App[Exec]
type DistExec Dist[Exec]
type BundleExec AppBundle[Exec]
type ProjectExec Project[Exec]

type ProjectGo interface {
	Project[Exec]
	IsGo()
}
