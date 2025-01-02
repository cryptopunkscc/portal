package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
)

type bundle struct {
	cacheDir string
	ctx      context.Context
	bundle   target.BundleExec
	cancel   context.CancelFunc
	args     []string
}

func Bundle(cacheDir string) target.Runner[target.BundleExec] {
	return &bundle{cacheDir: cacheDir}
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

	execFile, err := os.CreateTemp(r.cacheDir, r.bundle.Manifest().Package)
	if err != nil {
		return plog.Err(err)
	}

	e := r.bundle.Target().Executable()
	srcFile, err := e.Files().Open(e.Path())
	if err != nil {
		return plog.Err(err)
	}

	if err = execFile.Chmod(0755); err != nil {
		return plog.Err(err)
	}
	_, err = io.Copy(execFile, srcFile)
	if err != nil {
		return plog.Err(err)
	}
	if err = execFile.Close(); err != nil {
		return plog.Err(err)
	}
	defer os.Remove(execFile.Name())
	_ = srcFile.Close()

	var ctx context.Context
	ctx, r.cancel = context.WithCancel(r.ctx)
	err = Portal[target.Portal_](execFile.Name()).Run(ctx, r.bundle, r.args...)
	if err != nil {
		return err
	}
	return nil
}
