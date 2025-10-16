package portald

import (
	"errors"

	"github.com/cryptopunkscc/portal/pkg/resources"
)

type FsConfig struct {
	Repos map[string]string // list of paths to use for read-write storage
	Watch map[string]string // list of paths to use for read-only storage
}

func (s *Service) FsConfigRead() (c FsConfig, err error) {
	err = resources.Dir{s.Config.Astrald}.ReadYaml("fs.yaml", &c)
	return
}

func (s *Service) FsConfigWrite(c FsConfig) (err error) {
	return resources.Dir{s.Config.Astrald}.WriteYaml("fs.yaml", &c)
}

func (s *Service) FsConfigWatchAdd(name, path string) (err error) {
	c, err := s.FsConfigRead()
	if err != nil && !errors.Is(err, resources.ErrNotFound) {
		return
	}
	err = nil
	if c.Watch == nil {
		c.Watch = make(map[string]string)
	}
	c.Watch[name] = path
	return s.FsConfigWrite(c)
}
