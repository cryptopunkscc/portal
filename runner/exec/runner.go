package exec

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
	"os/exec"
)

type Runner struct{ portal.Config }

func DefaultRunner() (r Runner) {
	if err := r.Build(); err != nil {
		panic(err)
	}
	return
}

func WithReadWriter(ctx context.Context, rw io.ReadWriter) context.Context {
	return context.WithValue(ctx, stdKey, &Std{In: rw, Out: rw, Err: rw})
}

func (r Runner) RunApp(ctx context.Context, manifest target.Manifest, path string, args ...string) (err error) {
	tokens := token.Repository{Dir: r.Config.Tokens}
	t, err := tokens.Get(manifest.Package)
	if err != nil {
		return err
	}
	return r.Run(ctx, t.Token.String(), path, args...)
}

func (r Runner) Run(ctx context.Context, token string, path string, args ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("Command run: %s, %v", path, args)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = append(os.Environ(), apphost.AuthTokenEnv+"="+token)
	cmd.Env = append(cmd.Env, r.Env()...)
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
