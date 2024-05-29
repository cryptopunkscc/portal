package target

import (
	"fmt"
)

func DevPort(project Portal) string {
	return fmt.Sprintf("dev.%s", project.Manifest().Package)
}
