package pack

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/zip"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	"path"
)

func Run(_ context.Context, app target.Dist) (err error) {
	// create build dir
	buildDir := path.Join(app.Parent().Abs(), "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", app.Manifest().Name, app.Manifest().Version)
	if err = zip.Pack(app.Abs(), path.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("pack.Run: %v", err)
	}

	return
}
