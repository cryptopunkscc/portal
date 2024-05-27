package exec

import (
	"context"
	"os"
	osexec "os/exec"
)

type Service struct {
	executable string
}

func NewService(executable string) *Service {
	return &Service{executable: executable}
}

func (s Service) Run(ctx context.Context, cmd string, args ...string) error {
	args2 := append([]string{cmd}, args...)
	c := osexec.CommandContext(ctx, s.executable, args2...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
