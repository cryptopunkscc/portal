package target

import (
	"fmt"
)

func DevPort(project Project) string {
	return fmt.Sprintf("dev.%s", project.Manifest().Package)
}
