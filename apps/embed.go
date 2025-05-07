package apps

import (
	"embed"
)

//go:embed build
var Builds embed.FS

//go:embed launcher
var LauncherSvelteFS embed.FS

//go:embed profile
var ProfileFS embed.FS
