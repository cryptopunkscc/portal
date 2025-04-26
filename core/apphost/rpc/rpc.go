package rpc

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream/query"
)

type Rpc struct {
	Apphost apphost.Client
	Log     plog.Logger
	Codec   stream.Codec
}

func (r Rpc) codec() stream.Codec {
	c := r.Codec
	c.MarshalArgs = query.Marshal
	if c.Marshal == nil {
		c.Marshal = json.Marshal
		c.Ending = []byte{'\n'}
	}
	if c.Unmarshal == nil {
		c.Unmarshal = json.Unmarshal
	}
	return c
}
