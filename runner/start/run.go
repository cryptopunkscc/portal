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
	"sync"
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
	Open  bool   `cli:"open o"`
	Dev   bool   `cli:"dev d"`
}

func (s Start) Run(ctx context.Context, opt Opt, cmd ...string) (err error) {
	log := plog.Get(ctx).Type(s)
	s.portal.Logger(log)
	if err = s.portal.Ping(); err != nil {
		if err = startPortald(ctx, s.portal); err != nil {
			return
		}
	}
	wg := sync.WaitGroup{}
	if opt.Query != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(cmd) > 0 {
				time.Sleep(200 * time.Millisecond)
			}
			if ctx.Err() != nil {
				return
			}
			if err = s.queryApp(ctx, opt.Query); err != nil {
				return
			}
		}()
	}
	if len(cmd) > 0 {
		cmd = fixCmd(cmd)
		o := &portal.OpenOpt{}
		if opt.Dev {
			o.Schema = "dev"
			o.Order = []int{2, 1, 0}
		}
		if opt.Open {
			err = s.startApp(ctx, o, cmd)
		} else {
			err = s.runApp(ctx, o, cmd)
		}
	}
	log.Println("exit")
	wg.Wait()
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

func startPortaldProcess(ctx context.Context) (err error) {
	c := exec.Command("portald")
	err = c.Start()
	return
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

func fixCmd(args []string) []string {
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

func (s Start) startApp(ctx context.Context, opt *portal.OpenOpt, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting app:", cmd)
	return s.portal.Open(opt, cmd...)
}

func (s Start) runApp(ctx context.Context, opt *portal.OpenOpt, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("running app:", cmd)
	conn, err := s.portal.Connect(opt, cmd...)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		log.Println("reading done.")
	}()

	_, _ = io.Copy(os.Stdout, conn)
	log.Println("writing done.")
	return
}

func (s Start) queryApp(ctx context.Context, query string) (err error) {
	log := plog.Get(ctx)
	log.Println("running query:", query)

	conn, err := s.apphost.Query(id.Anyone, query)
	if err != nil {
		log.E().Printf("query (%s) FAILED: %v", query, err)
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
