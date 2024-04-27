package runtime

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
)

type New func() Api

type Api interface {
	apphost.Flat
}
