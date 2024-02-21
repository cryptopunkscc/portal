package main

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"log"
	"os"
)

func main() {
	js, err := apphost.GenerateAppHostJs(
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
	err = os.WriteFile("pkg/apphost/js/apphost.js", []byte(js), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
