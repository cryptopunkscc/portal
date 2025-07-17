package target

import (
	"github.com/cryptopunkscc/portal/api"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func init() {
	plog.Module = api.Module
}
