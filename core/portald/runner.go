package portald

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"sync"
)

type Runner[T Portal_] struct {
	cache     Cache[T]
	waitGroup sync.WaitGroup
	processes sig.Map[string, T]

	CacheDir string
	Shutdown func()
	Resolve  Resolve[T]
	Runners  func([]string) []Run[Portal_]
	Order    []int
}

func (s *Runner[T]) Run(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(s)
	log.Println("start")
	defer log.Println("exit")

	handler := cmd.Handler{
		Sub: s.Handlers(),
	}
	cmd.InjectHelp(&handler)
	router := apphost.Default.Rpc().Router(handler)
	router.Logger = log
	err = router.Run(ctx)

	if err != nil {
		log.Println(err)
		return nil
	}
	s.waitGroup.Wait()
	return
}
