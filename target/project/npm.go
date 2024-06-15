package project

import "github.com/cryptopunkscc/go-astral-js/target"

var _ target.ProjectNpm = (*nodeModule)(nil)

type nodeModule struct {
	portal
	target.NodeModule
}

func (m *nodeModule) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}
