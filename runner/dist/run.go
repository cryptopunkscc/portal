package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/npm"
	"os"
	"path"
)

func Run(m target.Project, dependencies []target.NodeModule) (err error) {
	if err = Prepare(m, dependencies); err != nil {
		return
	}
	if err = Dist(m); err != nil {
		return
	}
	return
}

func Prepare(m target.Project, dependencies []target.NodeModule) (err error) {
	if err = npm.Install(m); err != nil {
		return
	}
	if err = npm.InjectDependencies(m, dependencies); err != nil {
		return
	}
	return
}

func Dist(m target.Project) (err error) {
	if !m.CanNpmRunBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = npm.RunBuild(m); err != nil {
		return
	}
	if err = CopyIcon(m); err != nil {
		return
	}
	if err = CopyManifest(m); err != nil {
		return
	}
	return
}

func CopyIcon(m target.Project) (err error) {
	if m.Manifest().Icon == "" {
		return
	}
	iconSrc := path.Join(m.Abs(), m.Manifest().Icon)
	iconName := "icon" + path.Ext(m.Manifest().Icon)
	iconDst := path.Join(m.Abs(), "dist", iconName)
	if err = fs.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	m.Manifest().Icon = iconName
	return
}

func CopyManifest(m target.Project) (err error) {
	bytes, err := json.Marshal(m.Manifest())
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(m.Abs(), "dist", bundle.PortalJson), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
