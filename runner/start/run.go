package start

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	apphostRuntime "github.com/cryptopunkscc/portal/runtime/apphost"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func New(portal portal.Client, apphost apphost.Client) Start {
	return Start{portal: portal, apphost: apphost}
}

type Start struct {
	portal  portal.Client
	apphost apphost.Client
}

type Opt struct {
	Query string `cli:"query q"`
}

func (s Start) Run(ctx context.Context, opt Opt, args ...string) (err error) {
	plog.Get(ctx).Type(s).Set(&ctx)
	if err = s.portal.Ping(); err != nil {
		if err = startPortald(ctx, s.portal); err != nil {
			return
		}
	}
	if len(args) > 0 {
		args = fixArgs(args)
		if err = s.startApp(args); err != nil {
			return
		}
	}
	if opt.Query != "" {
		time.Sleep(200 * time.Millisecond) // Fixme
		if err = s.queryApp(ctx, opt.Query); err != nil {
			return
		}
	}
	if len(args) > 0 || opt.Query != "" {
		<-ctx.Done() // TODO exit automatically on portal or invoked app close
	}
	return
}

func startPortald(ctx context.Context, client portal.Client) (err error) {
	if err = startPortaldProcess(ctx); err != nil {
		return
	}
	if err = awaitPortaldService(ctx, client); err != nil {
		return
	}
	return
}

func startPortaldProcess(ctx context.Context) error {
	c := exec.CommandContext(ctx, "portald", "s")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}

func awaitPortaldService(ctx context.Context, client portal.Client) error {
	log := plog.Get(ctx)
	return flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
		if err = apphostRuntime.Init(); err != nil {
			return
		}
		return client.Ping()
	})
}

func fixArgs(args []string) []string {
	for i, arg := range args {
		args[i] = fixPath(arg)
	}
	return args
}

func fixPath(str string) string {
	if strings.HasPrefix(str, "./") || strings.HasPrefix(str, "../") {
		abs, err := filepath.Abs(str)
		if err != nil {
			panic(err)
		}
		return abs
	}
	return str
}

func (s Start) startApp(args []string) (err error) {
	return s.portal.Open(args...)
}

func (s Start) queryApp(ctx context.Context, query string) (err error) {
	log := plog.Get(ctx).Type(s)
	log.Println("Running query", query)

	conn, err := s.apphost.Query(id.Anyone, query)
	if err != nil {
		log.E().Printf("cannot query %s: %v", query, err)
		return err
	}

	defer conn.Close()
	c := make(chan any)
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		close(c)
	}()
	select {
	case <-ctx.Done():
	case <-c:
	}
	return
}
