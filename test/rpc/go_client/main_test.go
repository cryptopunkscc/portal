package main

import (
	"github.com/cryptopunkscc/go-astral-js/test/rpc"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	log.Println("\n=================================================================")
	time.Sleep(1 * time.Millisecond)
	rpc.NewTestClient("dev.test.%s").Run(t)
}
