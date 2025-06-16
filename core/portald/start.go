package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/core/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd/help"
)

func (s *Service[T]) Start(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(s)
	log.Println("starting portald...")
	ctx, s.shutdown = context.WithCancel(ctx)
	//s.Apphost.Log = log FIXME prints binary data & breaks tests
	if !s.configured {
		if err = s.Configure(); err != nil {
			return
		}
	}
	if err = s.startAstrald(ctx); err != nil {
		return
	}
	if err = s.startPortald(ctx); err != nil {
		return
	}
	s.createTokens(log)
	return
}

func (s *Service[T]) startAstrald(ctx context.Context) (err error) {
	r := astrald.Initializer{
		AgentAlias: "portald",
		NodeRoot:   s.Config.Astrald,
		TokensDir:  s.Config.Tokens,
		Config:     s.Config.Config,
		Runner:     s.Astrald,
		Apphost:    &s.Apphost,
	}
	return r.Start(ctx)
}

func (s *Service[T]) startPortald(ctx context.Context) error {
	log := plog.Get(ctx)
	handler := cmd.Handler{Sub: s.handlers()}
	help.Inject(&handler)
	router := s.Apphost.Rpc().Router(handler)
	if err := router.Init(ctx); err != nil {
		return err
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		if err := router.Listen(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("portald started")
	return nil
}

func (s *Service[T]) createTokens(log plog.Logger) {
	tokens := s.Tokens()
	for _, pkg := range s.ExtraTokens {
		if _, err := tokens.Resolve(pkg); err != nil {
			log.Println("cannot resolve token", err)
		}
	}
}
