package source

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io/fs"
)

type File_ struct {
	parent     Dir_
	pathToFile string
}

func (f File_) Abs() string                            { return f.parent.join(f.parent.Abs(), f.pathToFile) }
func (f File_) Path() string                           { return f.pathToFile }
func (f File_) File() (fs.File, error)                 { return f.FS().Open(f.pathToFile) }
func (f File_) FS() fs.FS                              { return f.parent.FS() }
func (f File_) Sub(_ ...string) (target.Source, error) { return nil, plog.Errorf("cannot sub file") }
func (f File_) IsDir() bool                            { return false }
