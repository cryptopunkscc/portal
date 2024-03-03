package binding

import (
	"embed"
)

//go:embed portal.js
var CommonJsString string

//go:embed portal.js
var CommonJsFs embed.FS
