//go:build !windows

package sig

import (
	"os"
	"syscall"
)

func Terminate() error {
	return syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func Interrupt() error {
	return syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
