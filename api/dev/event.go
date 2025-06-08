package dev

import (
	"time"
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

func NewMsg(pkg string, event Event) Msg {
	return Msg{Pkg: pkg, Event: event, Time: time.Now()}
}

type SendMsg func(msg Msg) error
