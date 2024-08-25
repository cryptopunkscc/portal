package query

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"strings"
)

type Command func(ctx context.Context) apphost.Client

func (cmd Command) Run(ctx context.Context, app string, args ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("Command", app, args)
	request := rpc.NewRequest(id.Anyone).Client(cmd(ctx))
	request.Logger(log)
	err = rpc.Command(request, app+" "+strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("command %s %v: %w", app, args, err)
	}
	return
}

func (cmd Command) Subscribe(ctx context.Context, app string, args ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("Command", app, args)
	request := rpc.NewRequest(id.Anyone).Client(cmd(ctx))
	request.Logger(log)
	err = rpc.Call(request, app+" "+strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("subscribe %s %v: %w", app, args, err)
	}
	scanner := bufio.NewScanner(request)
	for scanner.Scan() {
		println(scanner.Text())
	}
	return
}

func (cmd Command) Request(ctx context.Context, app string, args ...string) (err error) {
	if len(args) == 0 {
		_, err = Open.Run(ctx, app)
	} else {
		err = cmd.Subscribe(ctx, app, args...)
	}
	return
}
