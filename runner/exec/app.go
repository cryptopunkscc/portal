package exec

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
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
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.Env = append(os.Environ(), apphost.AuthTokenEnv+"="+token)
	cmd.Env = append(cmd.Env, d.Env()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err = tryRedirectStd(ctx, cmd); err != nil {
		return
	}
	if err = cmd.Run(); err != nil {
		err = plog.Err(err)
	}
	log.Printf("Command finished: %s, %v", path, args)
	return
}

func tryRedirectStd(ctx context.Context, cmd *exec.Cmd) (err error) {
	defer plog.TraceErr(&err)
	log := plog.Get(ctx)
	std, ok := ctx.Value(stdKey).(*Std)
	log.Printf("redirecting std %v ", ok)
	if !ok {
		return
	}
	cmd.SysProcAttr = appSysProcArgs()

	cmd.Stdout = nil
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	go io.Copy(std.Out, stdOut)

	cmd.Stderr = nil
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	go io.Copy(std.Err, stdErr)

	cmd.Stdin = nil
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	go func() {
		_, _ = io.Copy(stdIn, std.In)
		appKill(cmd)
		_ = stdIn.Close()
		log.Println("closed stdIn")
	}()
	return
}

const stdKey = "exec.RunCmd.Std"

type Std struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
