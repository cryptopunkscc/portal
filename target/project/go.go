package project

import "github.com/cryptopunkscc/go-astral-js/target"

var _ target.ProjectGo = (*golang)(nil)

type golang struct{ portal }

func (m *golang) DistGolang() target.DistExec {
	return Dist[target.DistExec](m)
}

func (m *golang) Type() target.Type {
	return target.TypeDev
}
