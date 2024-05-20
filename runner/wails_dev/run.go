package wails_dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	wailsdev "github.com/cryptopunkscc/go-astral-js/pkg/wails/dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
)

type Runner struct {
	frontCtx context.Context
	target.New
}

func NewRunner(bindings target.New) target.Run[target.ProjectFrontend] {
	return Runner{New: bindings}.Run
}

func (f Runner) Run(ctx context.Context, project target.ProjectFrontend) (err error) {
	log.Printf("portal dev open: (%d) %s\n", os.Getpid(), project.Manifest())
	defer log.Printf("portal dev close: (%d) %s\n", os.Getpid(), project.Manifest())
	opt := wails.AppOptions(f.New(target.TypeFrontend, "dev"))
	opt.OnStartup = func(ctx context.Context) {
		f.frontCtx = ctx
		go f.serve(project)
	}
	if err = wailsdev.Run(project.Abs(), opt); err != nil {
		log.Fatal(fmt.Errorf("dev.Run: %v", err))
	}
	return
}

func (f Runner) serve(project target.ProjectFrontend) {
	port := target.DevPort(project)
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), port+" ", 0))
	s.RouteFunc("reload", f.Reload)
	err := s.Run(f.frontCtx)
	if err != nil {
		log.Printf("%s: %v", port, err)
	}
}

func (f Runner) Reload() (err error) {
	if f.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(f.frontCtx)
	return
}
