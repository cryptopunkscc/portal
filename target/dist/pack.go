package dist

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/zip"
)

var PackRunner = target.SourceRunner[target.Dist_]{
	Resolve: Resolve_,
	Runner:  PackRun,
}

var PackRun target.Run[target.Dist_] = packRun

func packRun(_ context.Context, src target.Dist_, args ...string) (err error) {
	target.Op(&args, "clean")
	target.Op(&args, "pack")
	target.Op(&args, "goos=")
	target.Op(&args, "goarch=")
	return Pack(src, args...)
}

func Pack(app target.Dist_, path ...string) (err error) {
	if len(path) == 0 {
		path = []string{app.Abs(), ".."}
	}

	path = append(path, "build")
	buildDir := target.Abs(path...)

	if err = os.MkdirAll(buildDir, 0775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("os.MkdirAll: %v", err)
	}

	// pack dist dir
	platform := ""
	if d, ok := app.(target.DistExec); ok {
		t := d.Runtime().Target()
		if len(t.OS) > 0 {
			platform += "_" + t.OS
		}
		if len(t.Arch) > 0 {
			platform += "_" + t.Arch
		}
	}
	bundleName := fmt.Sprintf("%s_%s%s.portal", app.Manifest().Package, app.Version(), platform)
	if err = zip.Pack(app.Abs(), filepath.Join(buildDir, bundleName)); err != nil {
		return fmt.Errorf("pack.Run(%s): %v", app.Abs(), err)
	}
	return
}
