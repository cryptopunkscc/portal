package msg

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Broadcast struct {
	port    target.Port
	targets *sig.Map[string, target.Portal]
}

func NewBroadcast(
	port target.Port,
	targets *sig.Map[string, target.Portal],
) *Broadcast {
	return &Broadcast{
		port:    port,
		targets: targets,
	}
}

func (b *Broadcast) BroadcastMsg(ctx context.Context, msg target.Msg) {
	log := plog.Get(ctx).Type(b)
	log.Printf("msg %v, %v", msg, b.targets.Clone())
	for pkg := range b.targets.Clone() {
		if pkg == msg.Pkg {
			continue // skip sender
		}
		port := b.port.Copy(pkg)
		send := NewSend(port)
		if err := send(msg); err != nil {
			log.E().Println(err)
		}
	}
}
