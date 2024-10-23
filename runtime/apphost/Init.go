package apphost

import (
	"errors"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/mod/apphost/proto"
	"os"
	"strings"
)

// Based on github.com/cryptopunkscc/astrald/lib/astral/client.go

var defaultApphostAddrs = []string{"unix:~/.apphost.sock", "tcp:127.0.0.1:8625"}

// Check if astral.Client is successfully initialized.
func Check() (err error) {
	_, err = astral.GetNodeInfo(id.Anyone)
	return
}

// Init astral.Client with first available address and token.
func Init() error {
	if err := Check(); err == nil {
		return nil
	}

	token := os.Getenv(proto.EnvKeyToken)
	address := ""
	var addrs []string
	var envAddr = os.Getenv(proto.EnvKeyAddr)

	if len(envAddr) > 0 {
		addrs = strings.Split(envAddr, ";")
	} else {
		addrs = defaultApphostAddrs
	}

	for _, addr := range addrs {
		conn, err := proto.Dial(addr)
		if err == nil {
			conn.Close()
			address = addr
			break
		}
	}

	if address == "" {
		return ErrCannotConnect
	}

	astral.Client = *astral.NewClient(address, token)
	return nil
}

var ErrCannotConnect = errors.New("could not connect to apphost")
