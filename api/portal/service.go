package portal

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
)

type Service interface {
	Open() target.Request
	Shutdown() context.CancelFunc
}
