package golang

import (
	"github.com/cryptopunkscc/portal/api/target"
)

type Source struct {
	target.Project[target.Exec]
}

func (p *Source) IsGo()                       {}
func (p *Source) Changed(skip ...string) bool { return Changed(p, skip...) }
