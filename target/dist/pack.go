package dist

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/zip"
	"os"
	"path/filepath"
)

var PackRunner = target.SourceRunner[target.Dist_]{
	Resolve: Resolve_,
	Runner:  PackRun,
}

var PackRun target.Run[target.Dist_] = packRun

func packRun(_ context.Context, src target.Dist_, args ...string) (err error) {
	return Pack(src, args...)
}

func Pack(app target.Dist_, path ...string) (err error) {
	// create build dir
	//buildDir := filepath.Join(path...)
	//if buildDir == "" {
	//}
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
