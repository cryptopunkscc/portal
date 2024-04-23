package project

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"strings"
)

// Clones the given uri and returns the temporary cloned directory
func gitClone(uri string) (string, error) {
	// Create temporary directory
	dirname, err := os.MkdirTemp("", "wails-template-*")
	if err != nil {
		return "", err
	}

	// Parse remote template url and version number
	templateInfo := strings.Split(uri, "@")
	cloneOption := &git.CloneOptions{
		URL: templateInfo[0],
	}
	if len(templateInfo) > 1 {
		cloneOption.ReferenceName = plumbing.NewTagReferenceName(templateInfo[1])
	}

	_, err = git.PlainClone(dirname, false, cloneOption)

	return dirname, err
}
