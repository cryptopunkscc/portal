package fs2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Lock(path string, paths ...string) (unlock func(), err error) {
	path = filepath.Join(path, filepath.Join(paths...))
	lockFile, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil, errors.New("lock file already exists")
		}
		return nil, fmt.Errorf("failed to create lock file: %v\n", err)
	}
	return func() { _ = lockFile.Close() }, nil
}

func IsLocked(path string, paths ...string) bool {
	_, err := os.ReadFile(filepath.Join(path, filepath.Join(paths...)))
	return err == nil
}
