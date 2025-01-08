package sig

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"os/signal"
	"syscall"
)

func OnShutdown(log plog.Logger, cancel func()) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-quitChannel
	println()
	if log != nil {
		log.Println(sig, os.Args)
	}
	go cancel()
	<-quitChannel
	if log != nil {
		log.Println("force shutdown")
	}
	os.Exit(2)
}
