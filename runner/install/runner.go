package install

import (
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/mem"
)

type Runner struct {
	AppsDir mem.String
	Tokens  token.Repository
}
