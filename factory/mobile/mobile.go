package factory

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/feat/serve"
	runner "github.com/cryptopunkscc/portal/runner/mobile"
	runtime "github.com/cryptopunkscc/portal/runtime/mobile"
)

type mobile_ struct {
	api    mobile.Api
	ctx    context.Context
	cancel context.CancelFunc
	client apphost.Cached
}

func (m *mobile_) Api() mobile.Api      { return m.api }
func (m *mobile_) Ctx() context.Context { return m.ctx }
func (m *mobile_) Port() target.Port    { return target.PortPortal }
func (m *mobile_) Cache() apphost.Cache { return m.client }
func (m *mobile_) Serve() target.Request {
	s := &serve_{}
	s.Deps = s
	s.client = m.client
	s.astral = runtime.Astral{NodeRoot: m.api.NodeRoot()}.Run
	s.runHtmlApp = runner.Run[target.AppHtml](m.api.RequestHtml)
	s.CancelFunc = m.cancel
	return serve.Feat(s)
}
