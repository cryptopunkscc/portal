package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: mount,
		Name: "portal-drive",
		Desc: "Mounts astral objects fs to a directory.",
		Params: cmd.Params{
			{Type: "string", Desc: "Optional path to the directory to mount the fs to. Defaults to ~/portal-drive."},
		},
	}
}

func mount(ctx context.Context, dir string) (err error) {
	log := plog.Get(ctx)
	fs := NewFS(apphost.Default.Objects())
	defer fs.Close()

	if err = fs.Search(SearchArgs{
		Zone: astral.ZoneAll,
		Out:  "json",
	}); err != nil {
		return
	}

	drive := &GoFuseDrive{path: ".", fs: fs}

	if dir == "" {
		if dir, err = os.UserHomeDir(); err != nil {
			return
		}
		dir = filepath.Join(dir, "portal-drive")
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}

	if err = drive.Mount(dir); err != nil {
		return
	}

	go func() {
		<-ctx.Done()
		log.Printf("unmounting...")
		if err := drive.Unmount(); err != nil {
			log.Printf("unmount error: %v", err)
		}
	}()

	drive.Wait()
	return
}
