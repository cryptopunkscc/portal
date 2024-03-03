package binding

import (
	"embed"
)

//go:embed portal.js
var AllJsFS embed.FS

//go:embed portal.js
var AllJsString string
