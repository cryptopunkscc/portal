package create

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/create/templates"
	"github.com/flytam/filenamify"
	"log"
	"os"
	"os/exec"
)

func Run(
	projectName string,
	targetDir string,
	templateName string,
	force bool,
) (err error) {
	log.Println("force:", force)

	if targetDir == "" {
		targetDir = projectName
	}

	if force {
		if err = os.RemoveAll(targetDir); err != nil {
			return
		}
	}

	// prepare template options
	options := &templates.Options{
		ProjectName:  projectName,
		TargetDir:    targetDir,
		TemplateName: templateName,
	}
	if options.ProjectNameFilename, err = filenamify.Filenamify(
		options.ProjectName,
		filenamify.Options{Replacement: "_", MaxLength: 255},
	); err != nil {
		return
	}

	// generate project from template
	if _, _, err = templates.Install(options); err != nil {
		return
	}

	// npm install
	if err = npmInstall(targetDir); err != nil {
		return
	}

	// build project
	if err = build.Run(targetDir); err != nil {
		return
	}

	// bundle project
	if err = bundle.Run(targetDir); err != nil {
		return
	}

	return
}

func npmInstall(dir string) error {
	cmd := exec.Command("npm", "install")
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
