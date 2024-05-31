package broadcast

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"time"
)

type Msg struct {
	Pkg   string    `json:"pkg"`
	Event Event     `json:"signal"`
	Time  time.Time `json:"time"`
}

func NewMsg(pkg string, signal Event) Msg {
	return Msg{Pkg: pkg, Event: signal, Time: time.Now()}
}

type Event int

const (
	Changed Event = iota
	Refreshed
)

func Send(port string, msg Msg) error {
	request := rpc.NewRequest(id.Anyone, port)
	return rpc.Command(request, "", msg)
}
