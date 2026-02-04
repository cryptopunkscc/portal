package dev

import (
	"io"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
)

type Event int

const (
	Changed Event = iota
	Refreshed
)

type Msg struct {
	Pkg   string    `json:"pkg"`
	Event Event     `json:"event"`
	Time  time.Time `json:"time"`
}

var _ astral.Object = &Msg{}

func NewMsg(pkg string, event Event) *Msg {
	return &Msg{Pkg: pkg, Event: event, Time: time.Now()}
}

func (m Msg) ObjectType() string {
	return "portal.dev.msg"
}

func (m Msg) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Objectify(m).WriteTo(w)
}

func (m *Msg) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Objectify(m).ReadFrom(r)
}

type SendMsg func(msg *Msg) error
