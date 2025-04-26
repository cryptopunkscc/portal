package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
)

type tokenResponse struct{ i mod.AuthResponse }

func (t tokenResponse) Code() uint8               { return uint8(t.i.Code) }
func (t tokenResponse) GuestID() *astral.Identity { return t.i.GuestID }
func (t tokenResponse) HostID() *astral.Identity  { return t.i.HostID }
