package apps

import (
	"embed"
)

//go:embed */build */**/build
var Builds embed.FS

//go:embed launcher
var LauncherSvelteFS embed.FS
