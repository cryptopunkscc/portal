package target

import (
	"time"
)

type Event int

const (
	DevChanged Event = iota
	DevRefreshed
)

type Msg struct {
	Pkg   string    `json:"pkg"`
	Event Event     `json:"event"`
	Time  time.Time `json:"time"`
}

func NewMsg(pkg string, event Event) Msg {
	return Msg{Pkg: pkg, Event: event, Time: time.Now()}
}

type MsgSend func(msg Msg) error

type MsgSender func(port Port) MsgSend
