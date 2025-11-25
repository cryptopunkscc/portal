package dist

import (
	"encoding/json"
	"fmt"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
)

type Source_ struct {
	target.Source
	manifest manifest.Dist
}

var _ target.Dist_ = &Source_{}

func (s *Source_) IsApp()                       {}
func (s *Source_) IsDist()                      {}
func (s *Source_) Api() *manifest.Api           { return &s.manifest.Api }
func (s *Source_) Config() *manifest.Config     { return &s.manifest.Config }
func (s *Source_) Release() *manifest.Release   { return &s.manifest.Release }
func (s *Source_) Manifest() *manifest.App      { return &s.manifest.App }
func (s *Source_) MarshalJSON() ([]byte, error) { return json.Marshal(s.manifest) }
func (s *Source_) Version() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		s.manifest.Version, s.Api().Version, s.Release().Version,
	)
}
