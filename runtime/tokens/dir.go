package tokens

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"os"
	"path/filepath"
)

var Dir string
var Source target.Source

func init() {
	var err error
	if Dir, err = getDir(); err != nil {
		panic(err)
	}
	if Source, err = source.File(Dir); err != nil {
		panic(err)
	}
}

func getDir() (dir string, err error) {
	if dir, err = os.UserCacheDir(); err == nil {
		dir = filepath.Join(dir, "portal", "apps")
	} else {
		return
	}
	err = os.MkdirAll(dir, 0755)
	return
}
