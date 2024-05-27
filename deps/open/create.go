package open

import (
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/runner/app"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
)

func Create(
	portOpen string,
	wrapApi func(api target.Api) target.Api,
	findApps target.Find[target.App],
) target.Dispatch {
	runQuery := query.NewRunner[target.App](portOpen).Run
	newApphost := apphost.NewFactory(runQuery)
	newApi := target.ApiFactory(wrapApi,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
	runApp := app.NewRun(newApi)
	return open.NewFeat[target.App](findApps, runApp)
}
