package portal

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

type Client interface {
	Logger(logger plog.Logger)
	Join()
	Ping() (err error)
	Open(opt *OpenOpt, src ...string) error
	Connect(opt *OpenOpt, src ...string) (io.ReadWriteCloser, error)
	Close() error
}

type OpenOpt struct {
	Schema string `query:"s"`
	Order  []int  `query:"o"`
}
