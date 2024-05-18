package wails_dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	wailsdev "github.com/cryptopunkscc/go-astral-js/pkg/wails/dev"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
)

type Frontend struct {
	frontCtx context.Context
	target.New
	target.Project
}

func NewFrontend(bindings target.New, project target.Project) *Frontend {
	return &Frontend{New: bindings, Project: project}
}

func (f *Frontend) Start() (err error) {
	log.Printf("portal dev open: (%d) %s\n", os.Getpid(), f.Manifest())
	defer log.Printf("portal dev close: (%d) %s\n", os.Getpid(), f.Manifest())
	opt := wails.AppOptions(f.New(target.TypeFrontend, "dev"))
	opt.OnStartup = func(ctx context.Context) {
		f.frontCtx = ctx
		go f.serve()
	}
	if err = wailsdev.Run(f.Abs(), opt); err != nil {
		log.Fatal(fmt.Errorf("dev.Run: %v", err))
	}
	return
}

func (f *Frontend) serve() {
	port := target.DevPort(f.Project)
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), port+" ", 0))
	s.RouteFunc("reload", f.Reload)
	err := s.Run(f.frontCtx)
	if err != nil {
		log.Printf("%s: %v", port, err)
	}
}

func (f *Frontend) Reload() (err error) {
	if f.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(f.frontCtx)
	return
}
