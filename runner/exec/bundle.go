package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"io"
	"os"
)

type BundleRunner struct {
	cacheDir string
	ctx      context.Context
	bundle   target.BundleExec
	cancel   context.CancelFunc
}

func NewBundleRunner(cacheDir string) target.Runner[target.BundleExec] {
	return &BundleRunner{cacheDir: cacheDir}
}

func (r *BundleRunner) Run(ctx context.Context, bundle target.BundleExec) error {
	r.ctx = ctx
	r.bundle = bundle
	return r.Reload()
}

func (r *BundleRunner) Reload() error {
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
	err = NewPortal[target.Portal_](execFile.Name()).Run(ctx, r.bundle)
	if err != nil {
		return err
	}
	return nil
}
