package main

import (
	"github.com/cryptopunkscc/portal/core/portald/test/apps/go_client"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.SkipNow() // FIXME

	log.Println("\n=================================================================")
	time.Sleep(300 * time.Millisecond)
	go_client.NewTestClient(
		"test",
		"go",
		"js",
	).Run(t)
}
