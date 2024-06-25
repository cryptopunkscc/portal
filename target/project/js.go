package project

import "github.com/cryptopunkscc/portal/target"

var _ target.ProjectJs = (*js)(nil)

type js struct {
	target.ProjectNpm
	target.Js
}

func (m *js) DistJs() (t target.DistJs) {
	return Dist[target.DistJs](m)
}
