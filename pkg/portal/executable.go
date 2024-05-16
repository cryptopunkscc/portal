package portal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"time"
)

const Name = portal.Name

func Executable() string {
	executable, err := os.Executable()
	if err != nil {
		executable = "portal"
	}
	return executable
}

type CmdOpener[T target.Portal] struct {
	Resolve[T]
	action string
}

func NewCmdOpener[T target.Portal](resolve Resolve[T], action string) *CmdOpener[T] {
	return &CmdOpener[T]{Resolve: resolve, action: action}
}

func (o CmdOpener[T]) Open(ctx context.Context) func(src string, background bool) error {
	run := func(src string) error {
		// resolve apps from given source
		apps, err := o.Resolve(src)
		if len(apps) == 0 {
			return errors.Join(fmt.Errorf("CmdOpenerCtx %s %v no apps found in '%s'", o.action, reflect.TypeOf(o.Resolve), src), err)
		}
		// execute multiple targets as separate processes
		return Spawn[T](nil, ctx, apps, o.action)
	}

	return func(src string, background bool) error {
		if background {
			go run(src)
			return nil
		}
		return run(src)
	}
}

func Spawn[T target.Portal](
	wg *sync.WaitGroup,
	ctx context.Context,
	apps target.Portals[T],
	action string,
) (err error) {
	if wg == nil {
		wg = new(sync.WaitGroup)
	}
	wg.Add(len(apps))
	for _, t := range apps {
		log.Println(" * Spawn:", reflect.TypeOf(t), t.Manifest(), t.Abs())
		go func(p target.Portal) {
			defer wg.Done()
			if p.Type().Is(target.Frontend) {
				// TODO implement a better way to spawn backends before frontends
				time.Sleep(200 * time.Millisecond)
			}
			if err := CmdCtx(ctx, action, p.Abs(), "--attach").Run(); err != nil {
				return
			}
		}(t)
	}
	wg.Wait()
	return
}

func CmdCtx(ctx context.Context, args ...string) *exec.Cmd {
	var c *exec.Cmd
	if ctx != nil {
		c = exec.CommandContext(ctx, Executable(), args...)
	} else {
		c = exec.Command(Executable(), args...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}
