package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
)

func (a *Adapter) Fs() *FsClient {
	return &FsClient{a.Client}
}

type FsClient struct {
	*astrald.Client
}

func (client *FsClient) NewWatch(ctx *astral.Context, path, name string) (err error) {
	ch, err := client.QueryChannel(ctx, "fs.new_watch", query.Args{"path": path, "name": name})
	if err != nil {
		return
	}
	defer ch.Close()
	return ch.Switch(channel.ExpectAck, channel.PassErrors, channel.WithContext(ctx))
}
