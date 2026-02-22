package apphost

import (
	tree "github.com/cryptopunkscc/astrald/mod/tree/client"
)

func (a Adapter) Tree() *TreeClient { return &TreeClient{tree.New(a.TargetID, a.Client)} }

type TreeClient struct{ *tree.Client }
