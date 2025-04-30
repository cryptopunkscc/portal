package dist

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/api/target"
)

type Source_ struct {
	target.Source
	manifest *target.Manifest
}

var _ target.Dist_ = &Source_{}

func (s *Source_) IsApp()                       {}
func (s *Source_) IsDist()                      {}
func (s *Source_) Manifest() *target.Manifest   { return s.manifest }
func (s *Source_) MarshalJSON() ([]byte, error) { return json.Marshal(s.manifest) }
