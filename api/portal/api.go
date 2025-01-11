package portal

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

type Client interface {
	Logger(logger plog.Logger)
	Join()
	Ping() (err error)
	Open(src ...string) error
	Connect(src ...string) (io.ReadWriteCloser, error)
	Close() error
}
