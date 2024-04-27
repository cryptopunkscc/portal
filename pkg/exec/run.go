package exec

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func OnShutdown(cancel func()) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-quitChannel
	log.Println(sig)
	go cancel()
	<-quitChannel
	log.Println("force shutdown")
	os.Exit(2)
}

func Shutdown() error {
	return syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
