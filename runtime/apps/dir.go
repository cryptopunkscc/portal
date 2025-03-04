package apps

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"os"
	"path/filepath"
)

var Dir string
var Source target.Source

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	Dir = filepath.Join(dir, "portal", "apps")
	err = os.MkdirAll(Dir, 0755)
	if err != nil {
		panic(err)
	}
	if Source, err = source.File(Dir); err != nil {
		panic(err)
	}
}
