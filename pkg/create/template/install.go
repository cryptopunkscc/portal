package template

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Options for installing a template
type Options struct {
	Data
	Logger       *clilogger.CLILogger
	TemplateName string
	TargetDir    string
	InitGit      bool
}

func Install(opt *Options) (err error) {
	// Prepare target dir
	if opt.TargetDir != "" {
		// Try to install in opt.TargetDir
		if opt.TargetDir, err = filepath.Abs(opt.TargetDir); err != nil {
			return
		}
		if !fs.DirExists(opt.TargetDir) {
			if err = os.MkdirAll(opt.TargetDir, 0755); err != nil {
				return errors.Wrap(err, "fs.Mkdir")
			}
		}
	}

	// Resolve template
	template, err := getTemplateByShortname(opt.TemplateName)
	if err != nil {
		// Resolve absolute template path
		tempPath := opt.TemplateName
		if tempPath, err = filepath.Abs(opt.TemplateName); err != nil {
			return
		}

		switch {
		case fs.DirExists(tempPath):
			// Resolve from local dir
		default:
			// Resolve from git repository
			tempPath = ""
			defer func() {
				if err := os.RemoveAll(tempPath); err != nil {
					log.Fatal(err)
				}
			}()
			if tempPath, err = gitClone(opt.TemplateName); err != nil {
				return
			}

			// Remove the .git directory
			err = os.RemoveAll(filepath.Join(tempPath, ".git"))
			if err != nil {
				return
			}
		}
		// Parse template from path
		if template, err = parseTemplate(os.DirFS(tempPath)); err != nil {
			return errors.Wrap(err, "Error installing template")
		}
	}

	// Prepare gosod installer
	installer := gosod.New(template.FS)
	installer.IgnoreFile("template.json")
	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})

	// Prepare template data
	templateData := opt.Data
	if opt.AuthorName != "" {
		templateData.AuthorNameAndEmail = opt.AuthorName + " "
	}
	if opt.AuthorEmail != "" {
		templateData.AuthorNameAndEmail += "<" + opt.AuthorEmail + ">"
	}
	templateData.AuthorNameAndEmail = strings.TrimSpace(templateData.AuthorNameAndEmail)

	// Extract the template
	if err = installer.Extract(opt.TargetDir, templateData); err != nil {
		return errors.Wrap(err, "installer.Extract")
	}

	commonFs, err := debme.FS(templatesFs, "tmpl/common")
	if err != nil {
		return err
	}
	installer = gosod.New(commonFs)

	if err = installer.Extract(opt.TargetDir, templateData); err != nil {
		return errors.Wrap(err, "installer.Extract")
	}

	return
}

func InstallBase(src string) (err error) {
	commonFs, err := debme.FS(templatesFs, "tmpl/base")
	if err != nil {
		return err
	}
	installer := gosod.New(commonFs)

	if err = installer.Extract(src, struct{}{}); err != nil {
		return errors.Wrap(err, "installer.Extract")
	}
	return
}
