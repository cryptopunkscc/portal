package build

import (
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/list"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"os"
	"path"
)

func Run(dir string) (err error) {
	root := path.Clean(dir)
	dir = "."
	if !path.IsAbs(dir) {
		dir = root
		root, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	libs := list.Chan(project.Find[project.NodeModule](js.PortalLibFS, "."))
	if err = project.BuildPortalApps(root, dir, libs...); err != nil {
		return
	}
	return
}
