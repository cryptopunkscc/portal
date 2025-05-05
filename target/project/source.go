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

func (s *Source[T]) Api() *manifest.Api           { return &s.manifest.Api }
func (s *Source[T]) Build() *manifest.Builds      { return &s.manifest.Builds }
func (s *Source[T]) Config() *manifest.Config     { return &s.manifest.Config }
func (s *Source[T]) Manifest() *manifest.App      { return &s.manifest.App }
func (s *Source[T]) Runtime() T                   { return s.Dist().Runtime() }
func (s *Source[T]) Changed(skip ...string) bool  { return target.Changed(s, skip...) }
func (s *Source[T]) MarshalJSON() ([]byte, error) { return json2.Marshal(s.Manifest()) }
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
