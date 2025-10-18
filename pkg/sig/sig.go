package sig

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func OnShutdown(log plog.Logger, cancel func()) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-quitChannel
	println()
	if log != nil {
		log.Println(sig, os.Args)
	}
	cancel()
	<-quitChannel
	if log != nil {
		log.Println("force shutdown")
	}
	os.Exit(2)
}
