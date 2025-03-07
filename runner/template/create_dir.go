package template

import (
	"fmt"
	"os"
)

func checkTargetDirNotExist(dir, name string) (err error) {
	_, err = os.Stat(dir)
	if err == nil {
		err = fmt.Errorf("cannot create project %s: %s already exists", name, dir)
	} else if os.IsNotExist(err) {
		err = nil
	} else {
		err = fmt.Errorf("cannot create project %s: %v", name, err)
	}
	return
}

func makeTargetDir(dir string) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create project dir %s: %w", dir, err)
	}
	return
}
