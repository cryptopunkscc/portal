package dispatch

import (
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/runner/query"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func Create(executable string, queryOpen string) target.Dispatch {
	runQuery := query.NewRunner[target.App](queryOpen).Run
	return dispatch.NewFeat(executable, runQuery)
}
