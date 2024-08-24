package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/apps"
	createApphost "github.com/cryptopunkscc/portal/factory/apphost"
	createApps "github.com/cryptopunkscc/portal/factory/apps"
	"github.com/cryptopunkscc/portal/feat/apps"
	"log"
)

func main() {
	mod := module{}
	serve := apps.Feat(mod)
	err := serve(context.Background())
	if err != nil {
		log.Println(err)
	}
}

type module struct{}

func (m module) Apps() Apps             { return createApps.Default() }
func (m module) Client() apphost.Client { return createApphost.Client() }
