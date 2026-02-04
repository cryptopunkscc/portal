package main

import (
	"log"
	"testing"
	"time"

	"github.com/cryptopunkscc/portal/pkg/runner/deprecated/portald/test/apps/go_client"
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
