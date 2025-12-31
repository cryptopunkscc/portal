package initializer

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/keys"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (i *Astrald) readOrGenerateNodeIdentity() (err error) {
	if err = i.readNodeIdentity(); err != nil {
		err = i.generateNodeIdentity()
		return
	}
	return
}

func (i *Astrald) readNodeIdentity() (err error) {
	defer plog.TraceErr(&err)
	var pk keys.PrivateKey
	if err = i.resources.ReadObject("node_identity", &pk); err != nil {
		return
	}
	if i.nodeIdentity, err = astral.IdentityFromPrivKeyBytes(pk.Bytes); err != nil {
		return
	}
	i.log.Println("found existing node identity")
	return
}

func (i *Astrald) generateNodeIdentity() (err error) {
	defer plog.TraceErr(&err)
	i.nodeIdentity = astral.GenerateIdentity()
	if err = i.resources.WriteObject("node_identity", &keys.PrivateKey{
		Type:  keys.KeyTypeIdentity,
		Bytes: i.nodeIdentity.PrivateKey().Serialize(),
	}); err != nil {
		return
	}
	i.log.Println("generated node identity")
	return
}
