package apps

import (
	"embed"
)

//go:embed build
var Builds embed.FS

//go:embed launcher
var LauncherFS embed.FS

//go:embed profile
var ProfileFS embed.FS

//go:embed claim
var ClaimFS embed.FS
