package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"gopkg.in/yaml.v3"
)

type Dir struct {
	Path string
}

func (d Dir) Init() (err error) {
	defer plog.TraceErr(&err)
	if err = os.MkdirAll(d.Path, 0700); err != nil {
		return fmt.Errorf("cannot create resources directory: %w", err)
	}
	return
}

var ErrNotFound = errors.New("not found")

func (d Dir) Write(name string, data []byte) (err error) {
	defer plog.TraceErr(&err)
	if s, _ := os.Stat(d.Path); s == nil || !s.IsDir() {
		if err = d.Init(); err != nil {
			return
		}
	}
	return os.WriteFile(path.Join(d.Path, name), data, 0600)
}

func (d Dir) Read(name string) (bytes []byte, err error) {
	defer plog.TraceErr(&err)
	bytes, err = os.ReadFile(path.Join(d.Path, name))
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		err = ErrNotFound
	}
	return
}

type marshalFunc func(in interface{}) (out []byte, err error)

func (d Dir) Encode(marshal marshalFunc, name string, obj any) (err error) {
	bytes, err := marshal(obj)
	if err != nil {
		return
	}
	return d.Write(name, bytes)
}

type unmarshalFunc func(in []byte, out interface{}) (err error)

func (d Dir) Decode(unmarshal unmarshalFunc, name string, obj any) (err error) {
	bytes, err := d.Read(name)
	if err != nil {
		return
	}
	return unmarshal(bytes, obj)
}

func (d Dir) ReadYaml(name string, obj any) (err error)  { return d.Decode(yaml.Unmarshal, name, obj) }
func (d Dir) WriteYaml(name string, obj any) (err error) { return d.Encode(yaml.Marshal, name, obj) }
func (d Dir) ReadJson(name string, obj any) (err error)  { return d.Decode(json.Unmarshal, name, obj) }
func (d Dir) WriteJson(name string, obj any) (err error) { return d.Encode(json.Marshal, name, obj) }

func (d Dir) ReadObject(name string, obj astral.Object) (err error) {
	defer plog.TraceErr(&err)
	p := filepath.Join(d.Path, name)
	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	return ReadCanonical(file, obj)
}

func (d Dir) WriteObject(name string, obj astral.Object) (err error) {
	defer plog.TraceErr(&err)
	p := filepath.Join(d.Path, name)
	file, err := os.Create(p)
	if err != nil {
		return
	}
	defer file.Close()

	return WriteCanonical(file, obj)
}
