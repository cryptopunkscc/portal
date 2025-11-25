package rpc

import (
	"encoding/json"

	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream/object"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream/query"
)

type Rpc struct {
	Apphost apphost.Client
	Log     plog.Logger
	stream.Codec
}

func (r *Rpc) Format(name string) rpc.Rpc {
	switch name {
	case "json":
		r.Ending = []byte{'n'}
		r.Marshal = json.Marshal
		r.Unmarshal = json.Unmarshal
	case "object":
		r.Ending = nil
		r.Marshal = object.Marshal
		r.Unmarshal = object.Unmarshal
	default:
		plog.Printf("unknown rpc format %q", name)
	}
	return r
}

func (r *Rpc) codec() stream.Codec {
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
