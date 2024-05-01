package create

import (
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/template"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func Run(
	projectName string,
	targetDir string,
	templates []string,
	force bool,
) (err error) {
	log.Println("force:", force)

	if targetDir == "" {
		targetDir = projectName
	}
	if projectName == "" {
		// Get current working directory
		if projectName, err = os.Getwd(); err != nil {
			return
		}
		if targetDir == "" {
			targetDir = projectName
		}
		projectName = path.Base(projectName)
	}

	if force {
		if err = os.RemoveAll(targetDir); err != nil {
			return
		}
	}

	// prepare template options
	opt := template.Options{
		Data: template.Data{
			ProjectName: projectName,
		},
		TargetDir: targetDir,
	}

	// install each template
	for _, t := range templates {
		if len(t) == 0 {
			continue
		}
		c := strings.Split(t, ":")
		o := opt
		o.TemplateName = c[0]
		if len(c) > 1 {
			o.PackageName = c[1]
		} else {
			o.PackageName = c[0]
		}
		o.TargetDir = path.Join(o.TargetDir, o.PackageName)
		if err = runSingle(o); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func runSingle(opt template.Options) (err error) {
	// generate project from template
	if err = template.Install(&opt); err != nil {
		return errors.Wrap(err, "template.Install")
	}
	time.Sleep(500 * time.Millisecond)

	// npm install
	if err = npmInstall(opt.TargetDir); err != nil {
		return
	}
	time.Sleep(1000 * time.Millisecond)

	// build project
	if err = build.Run(opt.TargetDir); err != nil {
		return
	}
	time.Sleep(500 * time.Millisecond)

	// bundle project
	//if err = bundle.Run(opt.TargetDir); err != nil { TODO
	//	return
	//}
	time.Sleep(250 * time.Millisecond)

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
