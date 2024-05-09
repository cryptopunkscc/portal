package wails

import (
	"embed"
)

//go:embed portal.js
var JsString string

//go:embed portal.js
var JsFs embed.FS
