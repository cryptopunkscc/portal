package target2

import (
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

type Source interface {
	Abs() string
	Path() string
	Files() fs.FS
	Sub(dir string) (Source, error)
	IsDir() bool
}

type Template interface {
	Source
	Info() target.TemplateInfo
	Name() string
}

type Base interface {
	Source
	Manifest() *target.Manifest
}

type Portal[T any] interface {
	Base
	Target() T
}

type NodeModule interface {
	PkgJson() *target.PackageJson
}

type Project[T any] interface {
	Portal[T]
	IsProject()
	Dist() Dist[T]
}

type ProjectNpm[T any] interface {
	Project[T]
	NodeModule
}

type App[T any] interface {
	Portal[T]
	IsApp()
}

type Bundle interface {
	Source
	IsBundle()
}

type AppBundle[T any] interface {
	IsBundle()
	App[T]
}

type Dist[T any] interface {
	IsDist()
	App[T]
}

type Html interface{ IndexHtml() }
type PortalHtml Portal[Html]
type AppHtml App[Html]
type DistHtml Dist[Html]
type BundleHtml Bundle
type ProjectHtml ProjectNpm[Html]

type Js interface{ MainJs() }
type PortalJs Portal[Js]
type AppJs App[Js]
type DistJs interface{ Dist[Js] }
type BundleJs AppBundle[Js]
type ProjectJs ProjectNpm[Js]

type Exec interface{ Executable() Source }
type PortalExec Portal[Exec]
type AppExec App[Exec]
type DistExec Dist[Exec]
type BundleExec Bundle
type ProjectExec Project[Exec]

type ProjectGo interface {
	Project[Exec]
	IsGo()
}
