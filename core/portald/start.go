package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/core/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (s *Service[T]) Run(ctx context.Context) (err error) {
	if err = s.Start(ctx); err != nil {
		return
	}
	return s.Wait()
}

func (s *Service[T]) Start(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(s)
	log.Println("starting portald...")
	ctx, s.shutdown = context.WithCancel(ctx)
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
		NodeRoot:  s.NodeDir,
		TokensDir: s.TokensDir,
		Apphost:   &s.Apphost,
		Runner:    s.Astrald,
	}
	return r.Start(ctx)
}

func (s *Service[T]) startPortald(ctx context.Context) error {
	log := plog.Get(ctx)
	handler := cmd.Handler{Sub: s.handlers()}
	cmd.InjectHelp(&handler)
	router := s.Apphost.Rpc().Router(handler)
	router.Logger = log
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
	for _, pkg := range s.CreateTokens {
		if _, err := tokens.Resolve(pkg); err != nil {
			log.Println("cannot resolve token", err)
		}
	}
}
