package exec

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/exec"
)

func (r Runner) Bundle() *target.SourceRunner[target.BundleExec] {
	return &target.SourceRunner[target.BundleExec]{
		Resolve: target.Any[target.BundleExec](target.Try(exec.ResolveBundle)),
		Runner:  &BundleRunner{r},
	}
}

type BundleRunner struct{ Runner }

func (r *BundleRunner) Run(ctx context.Context, bundle target.BundleExec, args ...string) (err error) {
	execFile, err := r.unpackExecutable(bundle)
	if err != nil {
		return
	}

	err = r.RunApp(ctx, *bundle.Manifest(), execFile.Name(), args...)
	if err != nil {
		return
	}
	return
}

func (r *BundleRunner) unpackExecutable(bundle target.BundleExec) (execFile *os.File, err error) {
	defer plog.TraceErr(&err)

	src := bundle.Runtime().Executable()
	srcFile, err := src.File()
	defer srcFile.Close()
	if err != nil {
		return
	}
	srcId, err := readMD5Hex(srcFile)
	if err != nil {
		return
	}
	_ = srcFile.Close()
	if srcFile, err = src.File(); err != nil {
		return
	}

	execName := fmt.Sprintf("%s_%s_%s",
		bundle.Manifest().Package,
		bundle.Version(),
		srcId,
	)

	if execFile, err = os.OpenFile(filepath.Join(r.Bin, execName), os.O_RDWR|os.O_CREATE, 0755); err != nil {
		return
	}
	defer execFile.Close()
	execId, err := readMD5Hex(execFile)
	if err != nil {
		return
	}
	if err = execFile.Chmod(0755); err != nil {
		return
	}

	if execId == srcId {
		return
	}

	if _, err = io.Copy(execFile, srcFile); err != nil {
		return
	}
	return
}

func readMD5Hex(src io.Reader) (sum string, err error) {
	hash := md5.New()
	if _, err = io.Copy(hash, src); err != nil {
		err = plog.Err(err)
		return
	}
	bytes := hash.Sum(nil)
	sum = hex.EncodeToString(bytes)
	return
}
