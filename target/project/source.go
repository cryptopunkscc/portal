package project

import (
	json2 "encoding/json"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Source[T any] struct {
	target.Source
	manifest    manifest.Dev
	resolveDist target.Resolve[target.Dist[T]]
}

var _ target.Project[any] = &Source[any]{}

func (s *Source[T]) Api() *manifest.Api                        { return &s.manifest.Api }
func (s *Source[T]) Build() *manifest.Builds                   { return &s.manifest.Builds }
func (s *Source[T]) Config() *manifest.Config                  { return &s.manifest.Config }
func (s *Source[T]) Manifest() *manifest.App                   { return &s.manifest.App }
func (s *Source[T]) Changed(skip ...string) bool               { return target.Changed(s, skip...) }
func (s *Source[T]) MarshalJSON() ([]byte, error)              { return json2.Marshal(s.Manifest()) }
func (s *Source[T]) Dist_(platform ...string) (t target.Dist_) { return s.Dist(platform...) }
func (s *Source[T]) Dist(platform ...string) (t target.Dist[T]) {
	var err error
	defer plog.PrintTrace(&err)

	path := []string{"dist"}
	if len(platform) > 0 {
		path = append(path, platform...)
	}

	ss, err := s.Sub(path...)
	if err != nil {
		return
	}

	t, err = s.resolveDist(ss)
	return
}
func (s *Source[T]) Runtime() (t T) {
	d := s.Dist()
	if d != nil {
		t = d.Runtime()
	}
	return
}
