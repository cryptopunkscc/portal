package objects

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"gopkg.in/yaml.v3"
)

type FS struct {
	op      OpClient
	entries map[string]*Object
}

var _ fs.FS = &FS{}
var _ fs.ReadDirFS = &FS{}

func NewFS(client apphost.Client) *FS {
	return &FS{
		op:      Op(client.Rpc()),
		entries: map[string]*Object{},
	}
}

func (f *FS) Stat() (fs.FileInfo, error) { return f, nil }
func (f *FS) Read([]byte) (int, error)   { return 0, fs.ErrInvalid }
func (f *FS) Close() error               { return nil }
func (f *FS) Name() string               { return "." }
func (f *FS) Size() int64                { return 0 }
func (f *FS) Mode() fs.FileMode          { return fs.ModeDir | 0555 }
func (f *FS) ModTime() time.Time         { return time.Now() }
func (f *FS) IsDir() bool                { return true }
func (f *FS) Sys() any                   { return nil }

func (f *FS) Search(args SearchArgs) (err error) {
	search, err := f.op.Search(args)
	if err != nil {
		return
	}

	go func() {
		for r := range search {
			describe, err := f.op.Describe(DescribeArgs{
				ID:    *r.Object.ObjectID,
				Out:   "json",
				Zones: astral.ZoneAll,
			})
			if err != nil {
				plog.Println("cannot describe:", err)
				continue
			}
			file := &Object{
				ObjectID: *r.Object.ObjectID,
				Payload:  r.Payload,
				Describe: make(map[string]map[string]any),
			}
			for n := range describe {
				file.Describe[n["Type"].(string)] = n["Object"].(map[string]any)
			}
			if len(file.Payload) == 0 {
				file.Payload, _ = yaml.Marshal(file.Describe)
			}
			f.entries[file.Name()] = file
		}
	}()
	return
}

func (f *FS) ReadDir(name string) (entries []fs.DirEntry, err error) {
	for _, entry := range f.entries {
		entries = append(entries, entry)
	}
	return
}

func (f *FS) Open(name string) (file fs.File, err error) {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	if name == "" || name == "." {
		file = f
	} else if entry, ok := f.entries[name]; !ok {
		err = fs.ErrNotExist
	} else if file, err = os.Open(entry.Path()); err != nil {
		file = entry
		err = nil
	}
	return
}

type Object struct {
	astral.ObjectID
	Describe map[string]map[string]any
	Payload  []byte
}

var _ fs.File = &Object{}
var _ fs.FileInfo = &Object{}
var _ fs.DirEntry = &Object{}

func (f *Object) Read(bytes []byte) (n int, err error) { return copy(bytes, f.Payload), nil }
func (f *Object) Stat() (fs.FileInfo, error)           { return f, nil }
func (f *Object) Close() error                         { return nil }
func (f *Object) Name() string                         { return filepath.Base(f.Path()) }
func (f *Object) Size() int64                          { return int64(len(f.Payload)) }
func (f *Object) Mode() fs.FileMode                    { return 0400 }
func (f *Object) Type() fs.FileMode                    { return f.Mode().Type() }
func (f *Object) ModTime() time.Time                   { return time.UnixMicro(0) }
func (f *Object) IsDir() bool                          { return false }
func (f *Object) Sys() any                             { return nil }
func (f *Object) Info() (fs.FileInfo, error)           { return f, nil }
func (f *Object) Path() string {
	if f.Describe["mod.fs.file_location"] == nil {
		return f.ObjectID.String()
	}
	return f.Describe["mod.fs.file_location"]["Path"].(string)
}
