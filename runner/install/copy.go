package install

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func (i Runner) CopyOf(src target.Portal_) (err error) {
	return CopyFS(src.FS(), filepath.Join(i.AppsDir.Get(), src.Manifest().Package))
}

func CopyFS(source fs.FS, destPath string) error {
	return fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destPath, path)

		if d.IsDir() {
			if err = os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("mkdir %q: %v", targetPath, err)
			}
			return nil
		}
		srcFile, err := source.Open(path)
		if err != nil {
			return fmt.Errorf("open %q: %v", path, err)
		}
		defer srcFile.Close()

		//if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		//	return fmt.Errorf("mkdir %q: %v", targetPath, err)
		//}

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("create %q: %v", targetPath, err)
		}
		defer dstFile.Close()

		if _, err = io.Copy(dstFile, srcFile); err != nil {
			return fmt.Errorf("copy %q: %v", path, err)
		}

		//info, err := d.Info()
		//if err == nil {
		//	if err = os.Chmod(targetPath, info.Mode()); err != nil {
		//		return fmt.Errorf("chmod %q: %v", targetPath, err)
		//	}
		//}
		return nil
	})
}
