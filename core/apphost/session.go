package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type session struct{ i *lib.Session }

func (s session) Token(token string) (res api.TokenResponse, err error) {
	defer plog.TraceErr(&err)
	response, err := s.i.Token(token)
	if err != nil {
		return nil, err
	}
	return &tokenResponse{response}, nil
}

func (s session) Query(callerID *astral.Identity, targetID *astral.Identity, query string) (api.Conn, error) {
	return outConn(s.i.Query(callerID, targetID, query))
}

func (s session) Register(identity *astral.Identity, target string) (token string, err error) {
	return s.i.Register(identity, target)
}

func (s session) Close() error {
	return s.i.Close()
}
