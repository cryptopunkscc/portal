package core

import (
	"context"
	"encoding/json"

	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/html"
)

func (m *service) htmlRunner() *SourceRunner[AppHtml] {
	return &SourceRunner[AppHtml]{
		Resolve: Any[AppHtml](
			html.ResolveBundle.Try,
			html.ResolveDist.Try,
		),
		Runner: Run[AppHtml](m.runHtml),
	}
}

func (m *service) runHtml(_ context.Context, src AppHtml, args ...string) (err error) {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return
	}
	return m.mobile.StartHtml(src.Manifest().Package, string(argsJson))
}
