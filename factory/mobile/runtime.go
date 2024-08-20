package factory

import (
	"context"
	"github.com/cryptopunkscc/portal/api/mobile"
	feat "github.com/cryptopunkscc/portal/feat/mobile"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	runtime "github.com/cryptopunkscc/portal/runtime/mobile"
)

func Runtime(api mobile.Api) mobile.Runtime {
	client := apphost.Default()
	r := &runtime.Mobile{}
	r.Ctx, r.Cancel = context.WithCancel(context.Background())
	r.Client = client
	r.Serve = feat.Feat(&mobile_{
		ctx:    r.Ctx,
		cancel: r.Cancel,
		api:    api,
		client: client,
	})
	return r
}
