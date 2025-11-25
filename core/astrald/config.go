package astrald

import (
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (i *Initializer) createConfigs() (err error) {
	for name, config := range i.Config.Map() {
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
