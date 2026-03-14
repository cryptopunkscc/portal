package client

import (
	"context"
	"os"
	"path"
	"reflect"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/pkg/env"
	os2 "github.com/cryptopunkscc/portal/pkg/util/os"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

type Astrald struct {
	*astrald.Client
	Endpoint string
	Token    string
	Log      plog.Logger
	targetID *astral.Identity // Optional. HostID by default
}

func (a *Astrald) Init() *Astrald {
	a.Endpoint = firstNonZero[string](a.Endpoint, env.ApphostAddr.Get(), apphost.DefaultEndpoint)
	a.Token = firstNonZero[string](a.Token, env.ApphostToken.Get(), resolveTokenFromFile)
	a.Client = astrald.New(apphost.NewRouter(a.Endpoint, a.Token))
	return a
}

func firstNonZero[T any](items ...any) (v T) {
	for _, item := range items {
		switch x := item.(type) {
		case T:
			v = x
		case func() T:
			v = x()
		default:
			panic("invalid type")
		}
		if !reflect.ValueOf(v).IsZero() {
			return v
		}
	}
	return
}

func resolveTokenFromFile() string {
	abs := os2.Abs()
	file, err := os.Open(path.Join(abs, "astral_user"))
	if err != nil {
		return ""
	}
	defer file.Close()
	o, _, err := astral.Decode(file, astral.Canonical())
	if err != nil {
		return ""
	}
	return o.(*user.CreatedUserInfo).AccessToken.String()
}

func (a Astrald) WithTarget(identity *astral.Identity) *Astrald {
	if a.Client == nil {
		a.Init()
	}
	a.targetID = identity
	a.Client = a.Client.WithTarget(a.targetID)
	return &a
}

func (a Astrald) Resolve(name string) (i *astral.Identity, err error) {
	defer plog.TraceErr(&err)
	if name == "" || name == "localnode" {
		return a.HostID(), nil
	}
	if name == "self" {
		return a.GuestID(), nil
	}
	return a.Dir().ResolveIdentity(astral.NewContext(nil), name)
}

func (a *Astrald) Register(ctx context.Context) (out *astrald.Listener, err error) {
	defer plog.TraceErr(&err)
	listener, err := astrald.Listen()
	if err != nil {
		return nil, err
	}

	err = a.Apphost().RegisterHandler(astral.NewContext(ctx), listener.Endpoint(), listener.AuthToken())
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (a *Astrald) NodeAlias() (alias string, err error) {
	return a.Dir().GetAlias(nil, a.HostID())
}

func (a *Astrald) TargetID() *astral.Identity {
	return a.targetID
}
