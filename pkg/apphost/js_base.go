package apphost

import (
	"log"
)

var _jsBaseString string

func JsBaseString() string { return _jsBaseString }

func init() {
	var err error
	_jsBaseString, err = GenerateAppHostJs("builder", "default", "api", "jrpc")
	if err != nil {
		log.Fatalln(err)
	}
}
