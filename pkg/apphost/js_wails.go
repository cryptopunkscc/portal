package apphost

import "io/fs"

func JsWailsFs() fs.FS {
	sub, err := fs.Sub(jsFs, "js/wails")
	if err != nil {
		panic(err)
	}
	return sub
}
