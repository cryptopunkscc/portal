package runner

import (
	"archive/zip"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"io/fs"
	"os"
	"path"
	"strings"
)

type Target interface {
	Path() string
	Files() fs.FS
}

type GetTargets func(string) ([]Target, error)

func RawTargets(src string) (targets []Target, err error) {
	return findTargetDirs(src, targetProd)
}

func RawTargetsFS(files fs.FS) (targets []Target, err error) {
	return findTargetFS("", files, targetProd)
}

func DevTargets(src string) (targets []Target, err error) {
	return findTargetDirs(src, targetDev)
}

func ProdTargets(src string) (targets []Target, err error) {
	if targets, err = BundleTargets(src); len(targets) > 0 {
		return
	}
	if stat, err := os.Stat(src); err == nil && stat.IsDir() {
		if targets2, err := RawTargets(src); err == nil {
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
			targets = append(targets, project.NewDirectory(src, reader))
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
	return findTargetFS(dir, files, target)
}

func findTargetFS(dir string, files fs.FS, target fsTarget) (targets []Target, err error) {
	if err = fs.WalkDir(files, ".", func(p string, d fs.DirEntry, err error) error {
		if strings.Contains(p, "node_modules") {
			return fs.SkipDir
		}
		if sub := getTargetFS(files, p, d, target); sub != nil {
			t := project.NewDirectory(path.Join(dir, p), sub)
			targets = append(targets, t)
			return fs.SkipDir
		}
		return nil
	}); errors.Is(err, fs.SkipDir) {
		err = nil
	}
	return
}

func getTargetFS(files fs.FS, path string, d fs.DirEntry, target fsTarget) (project fs.FS) {
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
	if stat, err := fs.Stat(dir, bundle.PortalJson); err != nil || !stat.Mode().IsRegular() {
		return
	}
	if stat, err := fs.Stat(dir, "package.json"); err != nil || !stat.Mode().IsRegular() {
		return
	}
	return true
}

func targetProd(dir fs.FS) (b bool) {
	if stat, err := fs.Stat(dir, bundle.PortalJson); err != nil || !stat.Mode().IsRegular() {
		return
	}
	if stat, err := fs.Stat(dir, "package.json"); err == nil && stat.Mode().IsRegular() {
		return
	}
	return true
}
