package portal

import (
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"io"
)

var Request = rpc.NewRequest(id.Anyone, "portal")

func Bind(src string) (run func() error, closer io.ReadCloser, err error) {
	open, err := rpc.QueryFlow(id.Anyone, "portal.open")
	if err != nil {
		err = fmt.Errorf("portal.Bind failed: %v", err)
		return
	}
	run = func() error { return rpc.Command(open, "", src) }
	closer = open
	return
}
