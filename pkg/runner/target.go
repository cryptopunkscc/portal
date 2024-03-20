package runner

import (
	"archive/zip"
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"
)

type Target struct {
	Files fs.FS
	Path  string
}

type GetTargets func(string) ([]Target, error)

func DevTargets(src string) (targets []Target, err error) {
	return findTargetDirs(src, targetDev)
}

func ProdTargets(src string) (targets []Target, err error) {
	if targets, err = BundleTargets(src); len(targets) > 0 {
		return
	}
	if stat, err := os.Stat(src); err == nil && stat.IsDir() {
		if targets2, err := findTargetDirs(src, targetProd); err == nil {
			targets = append(targets, targets2...)
		}
	}
	return
}

func BundleTargets(src string) (targets []Target, err error) {
	if targets = appendZipTarget(targets, src); len(targets) > 0 {
		return
	}
	err = fs.WalkDir(os.DirFS(src), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(p, "node_modules") {
			return fs.SkipDir
		}
		targets = appendZipTarget(targets, path.Join(src, p))
		return nil
	})
	return
}

func appendZipTarget(targets []Target, src string) []Target {
	if path.Ext(src) == ".portal" {
		if reader, err := zip.OpenReader(src); err == nil {
			targets = append(targets, Target{
				Files: reader,
				Path:  src,
			})
		}
	}
	return targets
}

func findTargetDirs(dir string, target fsTarget) (targets []Target, err error) {
	_, err = os.Stat(dir)
	if err != nil {
		return
	}
	files := os.DirFS(dir)
	if err = fs.WalkDir(files, ".", func(p string, d fs.DirEntry, err error) error {
		if strings.Contains(p, "node_modules") {
			return fs.SkipDir
		}
		if sub := getTargetDir(files, p, d, target); sub != nil {
			t := Target{
				Files: sub,
				Path:  path.Join(dir, p),
			}
			targets = append(targets, t)
			return fs.SkipDir
		}
		return nil
	}); errors.Is(err, fs.SkipDir) {
		err = nil
	}
	return
}

func getTargetDir(files fs.FS, path string, d fs.DirEntry, target fsTarget) (project fs.FS) {
	if !d.IsDir() {
		return
	}
	sub, err := fs.Sub(files, path)
	if err != nil {
		return
	}
	if !target(sub) {
		return
	}
	project = sub
	return
}

type fsTarget func(fs.FS) bool

func targetDev(dir fs.FS) (b bool) {
	if stat, err := fs.Stat(dir, "portal.json"); err != nil || !stat.Mode().IsRegular() {
		return
	}
	if stat, err := fs.Stat(dir, "package.json"); err != nil || !stat.Mode().IsRegular() {
		return
	}
	return true
}

func targetProd(dir fs.FS) (b bool) {
	if stat, err := fs.Stat(dir, "portal.json"); err != nil || !stat.Mode().IsRegular() {
		return
	}
	if stat, err := fs.Stat(dir, "package.json"); err == nil && stat.Mode().IsRegular() {
		return
	}
	return true
}
