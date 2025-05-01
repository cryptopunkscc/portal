package template

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/leaanthony/gosod"
	"io/fs"
	"os"
	"path/filepath"
)

func (r *Factory) CreateProjects(targets map[string]string) (err error) {
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
		if err = r.createProject(t, n); err != nil {
			return
		}
	}
	return
}

func (r *Factory) createProject(template target.Template, projectName string) (err error) {
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

func (r *Factory) setupProjectArgs(projectName string) (args Args) {
	args = r.args
	args.ProjectName = projectName
	args.PackageName = "new.portal." + projectName
	return
}

func installTemplate(dir string, template fs.FS, args Args) (err error) {
	if err = extractTemplate(dir, template, args); err == nil {
		if err = extractCommons(dir, args); err == nil {
			return
		}
	}
	return fmt.Errorf("cannot install template: %w", err)
}

func extractTemplate(dir string, files fs.FS, args Args) (err error) {
	installer := gosod.New(files)
	installer.IgnoreFile("template.json")
	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("cannot extract template: %v", err)
	}
	return
}

func extractCommons(dir string, args Args) (err error) {
	installer := gosod.New(CommonsFs)
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("cannot extract template commons: %v", err)
	}
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
