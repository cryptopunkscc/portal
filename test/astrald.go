package test

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"strings"
)

type astrald struct {
	identity string
	alias    string
}

func (a *astrald) parseIdentity(logLine string) bool {
	c := strings.Split(logLine, " [node] astral node ")
	if len(c) < 2 {
		return false
	}
	a.identity = strings.TrimSpace(strings.Split(c[1], " ")[0])
	plog.Println("found node identity:", a.identity)
	return true
}

func (a *astrald) parseAlias(logLine string) bool {
	c := strings.Split(logLine, " [dir] call me ")
	if len(c) < 2 {
		return false
	}
	a.alias = strings.TrimSpace(c[1])
	plog.Println("found node alias:", a.alias)
	return true
}
