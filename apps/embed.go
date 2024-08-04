package apps

import (
	"embed"
)

//go:embed */build */**/build
var LauncherSvelteFS embed.FS
