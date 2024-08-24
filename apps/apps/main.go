package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/apps"
	create "github.com/cryptopunkscc/portal/factory/apps"
	"github.com/cryptopunkscc/portal/feat/apps"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
	"log"
)

func main() {
	mod := module{}
	serve := apps.Feat(mod)
	err := serve(context.Background())
	if err != nil {
		log.Println(err)
		log.Println(err)
	}
}

type module struct{}

func (m module) Apps() Apps             { return create.Default() }
func (m module) Client() apphost.Client { return runtime.Default() }
