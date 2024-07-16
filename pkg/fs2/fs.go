package fs2

import (
	"io"
	"os"
)

// CopyFile from source to target
func CopyFile(source string, target string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
