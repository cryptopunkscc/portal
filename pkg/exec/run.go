package exec

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Run(dir string, cmd ...string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
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
