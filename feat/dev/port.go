package dev

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func devPort(project target.Project) string {
	return fmt.Sprintf("dev.%s", project.Manifest().Package)
}
