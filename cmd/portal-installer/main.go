package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	println("Installing portal...\n")
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	bin := filepath.Join(home, ".local/bin")

	err = fs.WalkDir(binFs, "bin", func(srcPath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		dstPath := filepath.Join(bin, d.Name())
		print(fmt.Sprintf("* coping %s to %s", d.Name(), dstPath))

		dst, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
		if err != nil {
			return err
		}
		defer dst.Close()

		src, err := binFs.Open(srcPath)
		if err != nil {
			return err
		}
		defer src.Close()

		if _, err = io.Copy(dst, src); err != nil {
			_ = os.Remove(srcPath)
			return err
		}
		print(" [DONE]\n")
		return nil
	})
	if err != nil {
		panic(err)
	}
	println("\nPortal installed successfully")
}
