package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
)

type bundle struct {
	run    target.Run[target.BundleExec]
	ctx    context.Context
	cancel context.CancelFunc
	bundle target.BundleExec
	args   []string
}

func Bundle(cacheDir string) target.Runner[target.BundleExec] {
	return &bundle{
		run: BundleRun(cacheDir),
	}
}

func (r *bundle) Run(ctx context.Context, bundle target.BundleExec, args ...string) error {
	r.ctx = ctx
	r.bundle = bundle
	r.args = args
	return r.Reload()
}

func (r *bundle) Reload() error {
	if r.cancel != nil {
		r.cancel()
	}
	r.ctx, r.cancel = context.WithCancel(r.ctx)
	return r.run(r.ctx, r.bundle, r.args...)
}

func BundleRun(cacheDir string) target.Run[target.BundleExec] {
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
