package apphost

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed js/apphost.js
var _jsFs embed.FS

func JsFs() fs.FS {
	sub, err := fs.Sub(_jsFs, "js")
	if err != nil {
		panic(err)
	}
	return sub
}

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
		"static",
		"module",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
