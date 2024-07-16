package sig

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func OnShutdown(cancel func()) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-quitChannel
	log.Println(sig, os.Args)
	go cancel()
	<-quitChannel
	log.Println("force shutdown")
	os.Exit(2)
}

func Shutdown() error {
	return syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
