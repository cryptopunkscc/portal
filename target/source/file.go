package source

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type File_ struct {
	parent     Dir_
	pathToFile string
}

func (f File_) Abs() string                            { return f.parent.join(f.parent.Abs(), f.pathToFile) }
func (f File_) Path() string                           { return f.pathToFile }
func (f File_) File() (fs.File, error)                 { return f.FS().Open(f.pathToFile) }
func (f File_) FS() fs.FS                              { return f.parent.FS() }
func (f File_) Sub(_ ...string) (target.Source, error) { return nil, plog.Errorf("cannot sub file") }
func (f File_) IsDir() bool                            { return false }

func (f File_) CopyTo(path ...string) (err error) {
	defer plog.TraceErr(&err)

	srcFile, err := f.File()
	if err != nil {
		return
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return
	}

	destPath := filepath.Join(path...)
	destDir := filepath.Dir(destPath)
	if err = os.MkdirAll(destDir, 0755); err != nil {
		return
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return
	}
	if err = os.Chmod(destPath, srcInfo.Mode()); err != nil {
		return
	}
	return
}
