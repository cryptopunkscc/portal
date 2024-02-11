package apphost

import (
	"embed"
)

//go:embed js/apphost.js
var _jsFs embed.FS

func JsFs() embed.FS { return _jsFs }
