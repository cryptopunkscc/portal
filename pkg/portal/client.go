package portal

import (
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	jrpc "github.com/cryptopunkscc/go-apphost-jrpc"
	"io"
)

func Bind(src string) (run func() error, closer io.ReadCloser, err error) {
	open, err := jrpc.QueryFlow(id.Anyone, "portal.open")
	if err != nil {
		err = fmt.Errorf("portal.Bind failed: %v", err)
		return
	}
	run = func() error { return jrpc.Command(open, "", src) }
	closer = open
	return
}
