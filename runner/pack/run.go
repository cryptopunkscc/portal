package pack

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/zip"
	"os"
	"path/filepath"
)

func Run(_ context.Context, app target.Dist_) (err error) {
	// create build dir
	buildDir := filepath.Join(app.Abs(), "..", "build")
	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	bundleName := fmt.Sprintf("%s_%s.portal", app.Manifest().Package, app.Manifest().Version)
	if err = zip.Pack(app.Abs(), filepath.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("pack.Run: %v", err)
	}
	return
}
