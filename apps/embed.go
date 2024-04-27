package apps

import (
	"embed"
)

//go:embed launcher/svelte/dist
var LauncherSvelteFS embed.FS

//go:embed launcher/svelte/dist
var LauncherBackendFS embed.FS
