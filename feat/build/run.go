package build

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/arr"
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func Run(dir string) (err error) {
	libs := arr.FromChan(project.FindInFS[target.NodeModule](js.PortalLibFS))
	if err = project.BuildPortalApps(dir, ".", libs...); err != nil {
		return fmt.Errorf("cannot build portal apps: %w", err)
	}
	if err = project.BundlePortalApps(dir, "."); err != nil {
		return fmt.Errorf("cannot bundle portal apps: %w", err)
	}
	return
}
