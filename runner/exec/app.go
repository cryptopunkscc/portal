package exec

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
	"os/exec"
)

type RunApp func(ctx context.Context, manifest manifest.App, path string, args ...string) (err error)

type AppRunner struct{ portal.Config }

var _ RunApp = AppRunner{}.RunApp

func (d AppRunner) RunApp(ctx context.Context, manifest manifest.App, path string, args ...string) (err error) {
	tokens := token.Repository{Dir: d.Config.Tokens}
	t, err := tokens.Get(manifest.Package)
	if err != nil {
		return err
	}
	return d.run(ctx, t.Token.String(), path, args...)
}

func WithReadWriter(ctx context.Context, rw io.ReadWriter) context.Context {
	return context.WithValue(ctx, stdKey, &Std{In: rw, Out: rw, Err: rw})
}

func (d AppRunner) run(ctx context.Context, token string, path string, args ...string) (err error) {
	defer plog.TraceErr(&err)
	log := plog.Get(ctx).Type(d).Set(&ctx)
	log.Printf("Command run: %s, %v", path, args)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = append(os.Environ(), apphost.AuthTokenEnv+"="+token)
	cmd.Env = append(cmd.Env, d.Env()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err = setStd(cmd, ctx); err != nil {
		return
	}
	if err = cmd.Run(); err != nil {
		err = plog.Err(err)
	}
	log.Printf("Command finished: %s, %v", path, args)
	return
}

func setStd(cmd *exec.Cmd, ctx context.Context) (err error) {
	log := plog.Get(ctx)
	std, ok := ctx.Value(stdKey).(*Std)
	log.Printf("redirecting std %v ", ok)
	if !ok {
		return
	}

	cmd.Stdin = nil
	cmd.Stdout = std.Out
	cmd.Stderr = std.Err
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		err = plog.Err(err)
		return
	}
	go func() {
		_, _ = io.Copy(stdIn, std.In)
		_ = stdIn.Close()
	}()
	return
}

const stdKey = "exec.RunCmd.Std"

type Std struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
