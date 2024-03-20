package jrpc

import (
	"context"
	jrpc "github.com/cryptopunkscc/go-apphost-jrpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/create"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/dev"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/prod"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/publish"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
)

func Run(bindings runner.Bindings) (err error) {
	return jrpc.NewServer(
		func(context.Context, jrpc.Conn) *service {
			return newService(bindings)
		}).
		Logger(log.New(log.Writer(), "service ", 0)).
		Run(context.Background())
}

type service struct {
	bindings runner.Bindings
}

func (s service) String() string {
	return "portal"
}

func newService(bindings runner.Bindings) *service {
	return &service{bindings: bindings}
}

func (s *service) Create(
	projectName string,
	targetDir string,
	templates []string,
	force bool,
) (err error) {
	return create.Run(projectName, targetDir, templates, force)
}

func (s *service) Dev(src string) (err error) {
	return dev.Run(src, s.bindings)
}

func (s *service) Open(src string) (err error) {
	return prod.Run(src, s.bindings)
}

func (s *service) Build(src string) (err error) {
	return build.Run(src)
}

func (s *service) Bundle(src string) (err error) {
	return bundle.Run(src)
}

func (s *service) Publish(src string) (err error) {
	return publish.Run(src)
}
