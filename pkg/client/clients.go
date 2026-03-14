package client

import (
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/client"
	bip137sig "github.com/cryptopunkscc/astrald/mod/bip137sig/client"
	dir "github.com/cryptopunkscc/astrald/mod/dir/client"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	tree "github.com/cryptopunkscc/astrald/mod/tree/client"
	client "github.com/cryptopunkscc/portal/pkg/client/rpc"
	"github.com/cryptopunkscc/portal/pkg/util/rpc"
)

var Default = &Astrald{}

func init() {
	Default.Init()
}

func (a *Astrald) Apphost() *Apphost {
	return &Apphost{a.Client, *apphost.New(a.targetID, a.Client)}
}

func (a *Astrald) Dir() *dir.Client {
	return dir.New(a.targetID, a.Client)
}

func (a *Astrald) Tree() *tree.Client {
	return tree.New(a.targetID, a.Client)
}

func (a *Astrald) Objects() *Objects {
	return &Objects{*a.Client, *objects.New(a.targetID, a.Client)}
}

func (a *Astrald) Fs() *Fs {
	return &Fs{a.Client}
}

func (a *Astrald) Nodes() *Nodes {
	return &Nodes{*a.Client}
}

func (a *Astrald) User() *User {
	return &User{*a.Client}
}

func (a *Astrald) Rpc() rpc.Rpc {
	return &client.Rpc{Log: a.Log, Register: a.Register}
}

func (a *Astrald) Portald() Portald {
	return Portald{a}
}

func (a *Astrald) Bip127sig() *bip137sig.Client {
	return bip137sig.New(nil, a.Client)
}
