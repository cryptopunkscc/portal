package astrald

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/keys"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (r *Runner) readOrGenerateNodeIdentity() (err error) {
	if err = r.readNodeIdentity(); err != nil {
		err = r.generateNodeIdentity()
		return
	}
	return
}

func (r *Runner) readNodeIdentity() (err error) {
	defer plog.TraceErr(&err)
	var pk keys.PrivateKey
	if err = r.resources.DecodeObject("node_identity", &pk); err != nil {
		return
	}
	if r.nodeIdentity, err = astral.IdentityFromPrivKeyBytes(pk.Bytes); err != nil {
		return
	}
	r.log.Println("found existing node identity")
	return
}

func (r *Runner) generateNodeIdentity() (err error) {
	defer plog.TraceErr(&err)
	if r.nodeIdentity, err = astral.GenerateIdentity(); err != nil {
		return
	}
	if err = r.resources.EncodeObject("node_identity", &keys.PrivateKey{
		Type:  keys.KeyTypeIdentity,
		Bytes: r.nodeIdentity.PrivateKey().Serialize(),
	}); err != nil {
		return
	}
	r.log.Println("generated node identity")
	return
}
