package wails

import (
	"embed"
)

//go:embed portal.js
var WailsJsString string

//go:embed portal.js
var WailsJsFs embed.FS
