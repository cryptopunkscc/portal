package core

import (
	"context"
	"encoding/json"
	. "github.com/cryptopunkscc/portal/api/target"
)

func (m *service) htmlRun(_ context.Context, src App[Html], args ...string) (err error) {
	argsJson, err := json.Marshal(args)
	if err != nil {
		return
	}
	return m.mobile.StartHtml(src.Manifest().Package, string(argsJson))
}
