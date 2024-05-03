package portal

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"log"
	"time"
)

var name = "portal"

var Request = rpc.NewRequest(id.Anyone, name)

func init() {
	Request.Logger(log.New(log.Writer(), name+" ", 0))
}

func SrvOpen(src string) (pid int, err error) {
	pid, err = rpc.Query[int](Request, "open", src, true)
	return
}

func SrvOpenCtx(ctx context.Context, src string) (err error) {
	open, err := SrvOpenerCtx(ctx)
	if err != nil {
		err = fmt.Errorf("portal.open: %v", err)
		return
	}
	open(src)
	return
}

func SrvOpenerCtx(ctx context.Context) (open func(src string), err error) {
	var conn rpc.Conn
	if conn, err = rpc.QueryFlow(id.Anyone, "portal.open"); err != nil {
		err = fmt.Errorf("portal.opener: %v", err)
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

func Ping() (err error) {
	return rpc.Command(Request, "ping")
}

func Await(duration time.Duration) (err error) {
	err = Ping()
	_, err = exec.RetryT[any](duration, func(i int, n int, duration time.Duration) (_ any, err error) {
		err = Ping()
		return
	})
	return
}
