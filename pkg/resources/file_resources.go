package resources

import (
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/resources"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type FileResources struct {
	root string
	*resources.FileResources
}

func NewFileResources(root string, mkdir ...bool) (fr FileResources, err error) {
	plog.TraceErr(&err)
	fr.root = root
	fr.FileResources, err = resources.NewFileResources(root, len(mkdir) > 0 && mkdir[0])
	return
}

func (fr FileResources) ReadYaml(name string, obj any) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := fr.Read(name)
	if err != nil {
		return
	}
	return yaml.Unmarshal(bytes, obj)
}

func (fr FileResources) WriteYaml(name string, obj any) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := yaml.Marshal(obj)
	if err != nil {
		return
	}
	return fr.Write(name, bytes)
}

func (fr FileResources) DecodeObject(name string, o astral.Object) (err error) {
	defer plog.TraceErr(&err)
	path := filepath.Join(fr.root, name)
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	objType, payload, err := astral.OpenCanonical(file)
	switch {
	case err != nil:
		return
	case objType != o.ObjectType():
		return fmt.Errorf("invalid object type: %s", objType)
	}

	_, err = o.ReadFrom(payload)
	return
}

func (fr FileResources) EncodeObject(name string, o astral.Object) (err error) {
	defer plog.TraceErr(&err)
	path := filepath.Join(fr.root, name)
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = astral.WriteCanonical(file, o)
	return
}
