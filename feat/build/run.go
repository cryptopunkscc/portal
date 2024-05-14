package build

import (
	"github.com/cryptopunkscc/go-astral-js/feat/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/arr"
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
)

func Run(dir string) (err error) {
	base, sub, err := project.Path(dir)
	if err != nil {
		return
	}

	libs := arr.FromChan(project.Find[project.NodeModule](js.PortalLibFS, "."))
	if err = project.BuildPortalApps(base, sub, libs...); err != nil {
		return
	}
	return bundle.RunAll(dir)
}
