package portald

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/dir"
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
	ni.Identity, err = s.Apphost.Resolve("localnode")
	if err != nil {
		return
	}
	ni.Alias, err = dir.OpClient{Client: &s.Apphost}.GetAlias(*ni.Identity)
	if err != nil {
		return
	}
	s.NodeInfo = ni
	plog.Printf("NodeInfo: %s %s", s.Identity, s.Alias)
	return
}
