package build

import (
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/list"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
)

func Run(dir string) (err error) {
	libs := list.Chan(project.Find[project.NodeModule](js.PortalLibFS, "."))
	if err = project.BuildPortalApps(dir, ".", libs...); err != nil {
		return
	}
	return
}
