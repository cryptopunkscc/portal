package broadcast

import (
	"context"

	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/dev"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func New() *Service {
	return &Service{
		targets: &sig.Map[string, rpc.Conn]{},
	}
}

type Service struct {
	targets *sig.Map[string, rpc.Conn]
}

func (b *Service) BroadcastMsg(ctx context.Context, conn rpc.Conn) {
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
	b.targets.Replace(pkg, conn)

	// read messages
	var msg dev.Msg
	for {
		if msg, err = rpc.Decode[dev.Msg](conn); err != nil {
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
