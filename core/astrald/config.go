package astrald

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"path/filepath"
)

type Config struct {
	Node    core.Config    `yaml:",omitempty"`
	Apphost apphost.Config `yaml:",omitempty"`
}

func (i *Initializer) createConfigs() (err error) {
	for name, config := range map[string]any{
		"node":    i.Config.Node,
		"apphost": i.Config.Apphost,
	} {
		if err = i.writeIfNotExist(config, name+".yaml"); err != nil {
			return
		}
	}
	return
}

func (i *Initializer) writeIfNotExist(config any, name string) (err error) {
	defer plog.TraceErr(&err)
	abs := filepath.Join(i.NodeRoot, name)
	if _, err = os.Stat(abs); err != nil && !os.IsNotExist(err) {
		return
	}
	err = i.resources.WriteYaml(name, config)
	if err != nil {
		return
	}
	return
}
