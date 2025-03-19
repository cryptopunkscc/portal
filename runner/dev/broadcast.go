package dev

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/rpc"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func NewBroadcast() *Broadcast {
	return &Broadcast{
		targets: &sig.Map[string, rpc.Conn]{},
	}
}

type Broadcast struct {
	targets *sig.Map[string, rpc.Conn]
}

func (b *Broadcast) BroadcastMsg(ctx context.Context, conn rpc.Conn) {
	log := plog.Get(ctx).Type(b)

	// get caller package
	pkg, err := rpc.Decode[string](conn)
	if err != nil {
		return
	}
	// close previous if exist
	if v, ok := b.targets.Get(pkg); ok {
		_ = v.Close()
	}
	// append new
	b.targets.Set(pkg, conn)

	// read messages
	var msg target.Msg
	for {
		if msg, err = rpc.Decode[target.Msg](conn); err != nil {
			b.targets.Delete(pkg)
			return
		}
		log.Printf("msg %v, %v", msg, b.targets.Clone())

		for key, conn := range b.targets.Clone() {
			if key == msg.Pkg {
				continue // skip sender
			}
			if err := conn.Encode(msg); err != nil {
				log.E().Println(err.Error())
				b.targets.Delete(key)
			}
		}
	}
}
