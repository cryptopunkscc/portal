package main

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"log"
	"os"
	"strings"
)

func main() {
	request := apphost.Default.Rpc().Request("localnode", os.Args[1])
	request.Logger(plog.New().Scope("rpc"))
	var args []any = nil
	if len(os.Args) > 2 {
		args = append(args, strings.Join(os.Args[2:], " "))
	}
	c, err := rpc.Subscribe[any](request, "", args...)
	if err != nil {
		panic(err)
	}
	for a := range c {
		log.Println(a)
	}
}
