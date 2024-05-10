package portal

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"log"
	"strings"
	"time"
)

var name = "portal"

var Request = rpc.NewRequest(id.Anyone, name)

func init() {
	Request.Logger(log.New(log.Writer(), name+" ", 0))
}

func SrvOpenCtx(ctx context.Context, src ...string) (err error) {
	src = append([]string{}, src...)

	lastIndex := len(src) - 1
	action := src[lastIndex]
	prefix := src[:lastIndex]
	open, err := SrvOpenerCtx(prefix...)(ctx)
	if err != nil {
		return fmt.Errorf("portal.open: %v", err)
	}
	open(action)
	return
}

func SrvOpenerCtx(prefix ...string) func(ctx context.Context) (open func(src string), err error) {
	return func(ctx context.Context) (open func(src string), err error) {
		port := strings.Join(append(prefix, "portal.open"), ".")
		var conn rpc.Conn
		if conn, err = rpc.QueryFlow(id.Anyone, port); err != nil {
			err = fmt.Errorf("SrvOpenerCtx %s: %v", port, err)
			return
		}
		go func() {
			<-ctx.Done()
			_ = conn.Close()
		}()
		open = func(src string) {
			_ = rpc.Command(conn, "", src)
		}
		return
	}
}

func (r Runner[T]) Await(duration time.Duration) (err error) {
	err = r.Ping()
	_, err = exec.RetryT[any](duration, func(i int, n int, duration time.Duration) (_ any, err error) {
		err = r.Ping()
		return
	})
	return
}

func (r Runner[T]) Ping() (err error) {
	return rpc.Command(r.Request, "ping")
}
