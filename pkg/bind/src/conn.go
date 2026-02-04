package bind

import (
	"bufio"

	"github.com/cryptopunkscc/astrald/lib/apphost"
)

type conn struct {
	apphost.Conn
	bufio.Reader
}
