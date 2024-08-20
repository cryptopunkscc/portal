package mobile

import "github.com/cryptopunkscc/portal/api/bind"

type Runtime interface {
	Start()
	Stop()
	Apphost() Apphost
	Bindings(pkg string) bind.Runtime
}
