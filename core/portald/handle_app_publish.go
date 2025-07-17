package portald

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/source"
)

func (s *Service) PublishApps(path string) (out []bundle.Info, err error) {
	src, err := source.File(path)
	if err != nil {
		return
	}
	return s.PublishAppsFS(src)
}

func (s *Service) PublishAppsFS(src target.Source) (out []bundle.Info, err error) {
	p := s.Publisher()
	l := bundle.Resolve_.List(src)
	var id *astral.ObjectID
	var r *bundle.Release
	for _, b := range l {
		id, r, err = p.Publish(b)
		if err != nil {
			return
		}
		aa := bundle.Info{
			Manifest:  *b.Manifest(),
			Release:   *r,
			ReleaseID: id,
		}
		out = append(out, aa)
	}
	if len(l) == 0 {
		err = target.ErrNotFound
	}
	return
}
