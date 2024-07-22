package apps

import (
	"embed"
)

//go:embed launcher/svelte/build
var LauncherSvelteFS embed.FS
