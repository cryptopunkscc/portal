package test

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"strings"
)

type astrald struct {
	identity string
	alias    string
}

func (a *astrald) parseNodeInfo(logLine string) bool {
	c := strings.Split(logLine, ": NodeInfo: ")
	if len(c) < 2 {
		return false
	}
	c = strings.Split(c[1], " ")
	if len(c) < 2 {
		return false
	}
	a.identity = strings.TrimSpace(c[0])
	a.alias = strings.TrimSpace(c[1])
	plog.Println("found NodeInfo:", &a)
	return true
}
