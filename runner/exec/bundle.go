package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
)

func BundleRunner(cacheDir string) target.Run[target.BundleExec] {
	return func(ctx context.Context, bundle target.BundleExec, args ...string) (err error) {
		execFile, err := unpackExecutable(cacheDir, bundle)
		if err != nil {
			return
		}
		defer os.Remove(execFile.Name())

		err = RunCmd(ctx, execFile.Name(), args...)
		if err != nil {
			return
		}
		return
	}
}

func unpackExecutable(cacheDir string, bundle target.BundleExec) (execFile *os.File, err error) {
	execName := fmt.Sprintf("%s_%s", bundle.Manifest().Package, bundle.Manifest().Version)
	execFile, err = os.CreateTemp(cacheDir, execName)
	//execFile, err = os.Create(filepath.Join(cacheDir, execName)) // FIXME
	if err != nil {
		return nil, plog.Err(err)
	}

	execSource := bundle.Target().Executable()
	execSrcFile, err := execSource.Files().Open(execSource.Path())
	if err != nil {
		return nil, plog.Err(err)
	}

	if err = execFile.Chmod(0755); err != nil {
		return nil, plog.Err(err)
	}
	_, err = io.Copy(execFile, execSrcFile)
	if err != nil {
		return nil, plog.Err(err)
	}

	if err = execFile.Close(); err != nil {
		return nil, plog.Err(err)
	}
	_ = execSrcFile.Close()
	return
}
