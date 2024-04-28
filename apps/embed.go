package apps

import (
	"embed"
)

//go:embed launcher/svelte/dist
var LauncherSvelteFS embed.FS
