package install

import (
	"github.com/cryptopunkscc/portal/core/token"
)

type Runner struct {
	AppsDir string
	Tokens  token.Repository
}
