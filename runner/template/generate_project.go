package template

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/template"
	"io/fs"
	"os"
	"path/filepath"
)

func (r *Runner) GenerateProjects(targets map[string]string) (err error) {
	for _, t := range List() {
		if n, ok := targets[t.Name()]; ok {
			if err = r.GenerateProject(t, n); err != nil {
				return
			}
		}
	}
	return
}

func (r *Runner) GenerateProject(template target.Template, projectName string) (err error) {
	dir := filepath.Join(r.dir, projectName)

	if err = checkTargetDirNotExist(dir, projectName); err != nil {
		return
	}
	if err = makeTargetDir(dir); err != nil {
		return
	}

	args := r.setupProjectArgs(projectName)
	if err = installTemplate(dir, template.Files(), args); err != nil {
		return
	}

	if err = fixupManifestIfNeeded(template, dir); err != nil {
		return
	}
	return
}

func (r *Runner) setupProjectArgs(projectName string) (args template.Args) {
	args = r.args
	args.ProjectName = projectName
	args.PackageName = "new.portal." + projectName
	return
}

func fixupManifestIfNeeded(t target.Template, dir string) (err error) {
	stat, err := fs.Stat(t.Files(), "dev")
	shouldFix := err == nil && !stat.IsDir()
	if !shouldFix {
		return
	}
	if err = os.Rename(
		filepath.Join(dir, "portal.json"),
		filepath.Join(dir, "dev.portal.json"),
	); err != nil {
		err = fmt.Errorf("cannot fix up manifest: %w", err)
	}
	return
}
