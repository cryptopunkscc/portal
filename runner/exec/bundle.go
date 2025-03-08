package exec

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/tokens"
	"io"
	"os"
	"path/filepath"
)

func BundleRunner(cacheDir string) target.Run[target.BundleExec] {
	return func(ctx context.Context, bundle target.BundleExec, args ...string) (err error) {
		execFile, err := unpackExecutable(cacheDir, bundle)
		if err != nil {
			return
		}

		token, err := tokens.Repository{}.Get(bundle.Manifest().Package)
		if err != nil {
			return err
		}

		err = RunCmd(ctx, token.Token.String(), execFile.Name(), args...)
		if err != nil {
			return
		}
		return
	}
}

func HostBundleRunner(cacheDir string, token string) target.Run[target.BundleExec] {
	return func(ctx context.Context, bundle target.BundleExec, args ...string) (err error) {
		execFile, err := unpackExecutable(cacheDir, bundle)
		if err != nil {
			return
		}
		defer os.Remove(execFile.Name())

		err = RunCmd(ctx, token, execFile.Name(), args...)
		if err != nil {
			return
		}
		return
	}
}

func unpackExecutable(cacheDir string, bundle target.BundleExec) (execFile *os.File, err error) {
	binDir := filepath.Join(cacheDir, "bin")
	if err = os.MkdirAll(binDir, 0755); err != nil {
		return nil, plog.Err(err)
	}

	src := bundle.Target().Executable()
	srcFile, err := src.File()
	defer srcFile.Close()
	if err != nil {
		return nil, plog.Err(err)
	}
	srcId, err := readMD5Hex(srcFile)
	if err != nil {
		return nil, plog.Err(err)
	}
	_ = srcFile.Close()
	srcFile, err = src.File()
	if err != nil {
		return nil, plog.Err(err)
	}

	execName := fmt.Sprintf("%s_%s_%s",
		bundle.Manifest().Package,
		bundle.Manifest().Version,
		srcId,
	)

	execFile, err = os.OpenFile(filepath.Join(binDir, execName), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, plog.Err(err)
	}
	defer execFile.Close()
	execId, err := readMD5Hex(execFile)
	if err != nil {
		return nil, plog.Err(err)
	}
	if err = execFile.Chmod(0755); err != nil {
		return nil, plog.Err(err)
	}

	if execId == srcId {
		return execFile, nil
	}

	if _, err = io.Copy(execFile, srcFile); err != nil {
		return nil, plog.Err(err)
	}
	return
}

func readMD5Hex(src io.Reader) (sum string, err error) {
	hash := md5.New()
	if _, err = io.Copy(hash, src); err != nil {
		return
	}
	bytes := hash.Sum(nil)
	sum = hex.EncodeToString(bytes)
	return
}
