package project

import (
	json2 "encoding/json"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
)

func New[T any](
	source target.Source,
	resolve target.Resolve[T],
) (project target.Project[T], err error) {
	defer plog.TraceErr(&err)
	if _, err = resolve(source); err != nil {
		return
	}
	s := &Source[T]{}
	if err = all.Unmarshalers.Load(&s.manifest, source.FS(), target.BuildFilename); err != nil {
		return
	}
	s.build = target.LoadBuilds(source)
	s.resolveDist = dist.Resolver(resolve)
	s.Source = source
	if s.manifest.Exec == "" {
		s.manifest.Exec = target.GetBuild(s).Exec
	}
	project = s
	return
}

type Source[T any] struct {
	target.Source
	build       target.Builds
	manifest    target.Manifest
	resolveDist target.Resolve[target.Dist[T]]
}

func (s *Source[T]) Changed(skip ...string) bool  { return target.Changed(s, skip...) }
func (s *Source[T]) MarshalJSON() ([]byte, error) { return json2.Marshal(s.Manifest()) }
func (s *Source[T]) Manifest() *target.Manifest   { return &s.manifest }
func (s *Source[T]) Target() T                    { return s.Dist().Target() }
func (s *Source[T]) Build() target.Builds         { return s.build }
func (s *Source[T]) Dist_() (t target.Dist_)      { return s.Dist() }
func (s *Source[T]) Dist() (t target.Dist[T]) {
	sub, err := s.Sub("dist")
	if err != nil {
		plog.Println(plog.Err(err))
		return
	}
	t, err = s.resolveDist(sub)
	if err != nil {
		plog.Println(plog.Err(err))
		return
	}
	return
}
