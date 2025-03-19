package template

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"io/fs"
	"os"
	"path/filepath"
)

func (r *Runner) GenerateProjects(targets map[string]string) (err error) {
	availableTemplates := Map()
	matched := map[string]target.Template{}
	for projectName, templateName := range targets {
		t, ok := availableTemplates[templateName]
		if !ok {
			return fmt.Errorf("project template %s not found", templateName)
		}
		matched[projectName] = t
	}
	for n, t := range matched {
		if err = r.GenerateProject(t, n); err != nil {
			return
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
	if err = installTemplate(dir, template.FS(), args); err != nil {
		return
	}

	if err = fixupManifestIfNeeded(template, dir); err != nil {
		return
	}
	return
}

func (r *Runner) setupProjectArgs(projectName string) (args Args) {
	args = r.args
	args.ProjectName = projectName
	args.PackageName = "new.portal." + projectName
	return
}

func fixupManifestIfNeeded(t target.Template, dir string) (err error) {
	stat, err := fs.Stat(t.FS(), "dev")
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
