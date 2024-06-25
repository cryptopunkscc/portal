package exec

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	targetSource "github.com/cryptopunkscc/portal/target/source"
	"io/fs"
)

func Resolve(portal target.Portal) (t target.Exec, err error) {
	file := portal.Manifest().Exec
	stat, err := fs.Stat(portal.Files(), file)
	if err != nil {
		return
	}
	if stat.Mode().Perm()&0111 == 0 {
		err = plog.Errorf("not executable %s", file)
		return
	}
	executable := targetSource.FromFS(portal.Files(), file, portal.Abs()).Lift()
	t = &source{executable: executable}
	return
}
