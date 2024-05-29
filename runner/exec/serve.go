package exec

import (
	"context"
	"os"
	"os/exec"
)

type Service struct {
	executable string
}

func NewService(executable string) *Service {
	return &Service{executable: executable}
}

func (s Service) Start(ctx context.Context, cmd string, args ...string) error {
	args2 := append([]string{cmd}, args...)
	c := exec.CommandContext(ctx, s.executable, args2...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
