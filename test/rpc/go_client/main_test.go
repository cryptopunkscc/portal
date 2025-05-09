package main

import (
	"github.com/cryptopunkscc/portal/test/rpc"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.SkipNow() // FIXME

	log.Println("\n=================================================================")
	time.Sleep(300 * time.Millisecond)
	rpc.NewTestClient(
		"test",
		"go",
		"js",
	).Run(t)
}
