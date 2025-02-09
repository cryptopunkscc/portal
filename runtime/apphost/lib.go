package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

var Lib = &apphost.Client{}

func IsConnected() bool {
	return Lib.HostID != nil
}

func Connect(ctx context.Context) (err error) {
	if IsConnected() {
		return
	}
	return Reconnect(ctx)
}

func Reconnect(ctx context.Context) (err error) {
	l, err := apphost.NewDefaultClient()
	if err == nil {
		Lib.HostID = l.HostID
		Lib.GuestID = l.GuestID
		Lib.AuthToken = l.AuthToken
		Lib.Endpoint = l.Endpoint
	}
	log := plog.Get(ctx)
	log.Printf("host id: %s", Lib.HostID.String())
	log.Printf("guest id: %s", Lib.GuestID.String())
	log.Printf("token: %s", Lib.AuthToken)
	log.Printf("endpoint: %s", Lib.Endpoint)
	return
}
