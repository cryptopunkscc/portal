package factory

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/multi"
)

type serve_ struct {
	srv.Module[target.App_]
	astral     serve.Astral
	runHtmlApp target.Run[target.AppHtml]
	client     apphost.Client
}

func (s *serve_) Astral() serve.Astral                 { return s.astral }
func (s *serve_) Client() apphost.Client               { return s.client }
func (s *serve_) Resolve() target.Resolve[target.App_] { return apps.ResolveAll }

func (s *serve_) Run() target.Run[target.App_] {
	return multi.Runner[target.App_](
		app.Run(goja.NewRun(bind.BackendRuntime())),
		app.Run(s.runHtmlApp),
	)
}

func (s *serve_) Priority() target.Priority {
	return []target.Matcher{
		target.Match[target.Bundle_],
		target.Match[target.Dist_],
	}
}
