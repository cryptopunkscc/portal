package exec

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
	"path/filepath"
)

func BundleRunner() target.Run[target.BundleExec] {
	return func(ctx context.Context, bundle target.BundleExec, args ...string) (err error) {
		execFile, err := unpackExecutable(bundle)
		if err != nil {
			return
		}

		err = Cmd{}.RunApp(ctx, *bundle.Manifest(), execFile.Name(), args...)
		if err != nil {
			return
		}
		return
	}
}

func unpackExecutable(bundle target.BundleExec) (execFile *os.File, err error) {
	defer plog.TraceErr(&err)
	binDir := env.PortaldBin.MkdirAll()
	if len(binDir) == 0 {
		return nil, plog.Errorf("no executable path specified")
	}

	src := bundle.Target().Executable()
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
		bundle.Manifest().Version,
		srcId,
	)

	if execFile, err = os.OpenFile(filepath.Join(binDir, execName), os.O_RDWR|os.O_CREATE, 0755); err != nil {
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
