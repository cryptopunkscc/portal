package apps

import (
	"embed"
)

//go:embed */build */**/build
var FS embed.FS
