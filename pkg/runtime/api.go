package runtime

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type New func(p target.Type) Api

type Api interface {
	apphost.Flat
}
