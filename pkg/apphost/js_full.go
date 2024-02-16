package apphost

import (
	"embed"
	"log"
)

//go:embed js/apphost.js
var _jsFs embed.FS

func JsFs() embed.FS { return _jsFs }

var _jsString string

func JsString() string { return _jsString }

func JsBytes() []byte { return []byte(_jsString) }

func init() {
	var err error
	_jsString, err = GenerateAppHostJs(
		"builder",
		"wails",
		"android",
		"default",
		"api",
		"jrpc",
		"module",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
