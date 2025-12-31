package portald

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type NodeInfo struct {
	Identity *astral.Identity `json:"identity"`
	Alias    string           `json:"alias"`
}

func (s *Service) nodeInfo() (ni *NodeInfo, err error) {
	if s.NodeInfo != nil {
		return s.NodeInfo, nil
	}
	ni = new(NodeInfo)
	ni.Alias, err = s.Apphost.NodeAlias()
	if err != nil {
		return
	}
	ni.Identity = s.Apphost.HostID()
	s.NodeInfo = ni
	plog.Printf("NodeInfo: %s %s", s.Identity, s.Alias)
	return
}
