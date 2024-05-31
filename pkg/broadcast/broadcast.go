package broadcast

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"strings"
)

type DevBroadcast struct {
	prefix  []string
	command string
	targets *sig.Map[string, target.Portal]
}

func New(
	targets *sig.Map[string, target.Portal],
	command string,
	prefix ...string,
) *DevBroadcast {
	return &DevBroadcast{
		targets: targets,
		command: command,
		prefix:  prefix,
	}
}

func (b *DevBroadcast) Signal(ctx context.Context, msg Msg) {
	log := plog.Get(ctx).Type(b)
	log.Printf("msg %v, %v", msg, b.targets.Clone())
	for p := range b.targets.Clone() {
		if p == msg.Pkg {
			continue
		}
		port := strings.Join(append(b.prefix, p, "ctrl"), ".")
		if err := Send(port, msg); err != nil {
			log.E().Println(err)
		}
	}
}
