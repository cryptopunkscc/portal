package core

import (
	"strings"

	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/source/html"
)

type htmlBundleRunner struct {
	html.App
	api    mobile.Api
	Bundle html.Bundle
}

func (r htmlBundleRunner) New() source.Source {
	return &r
}

func (r *htmlBundleRunner) ReadSrc(src source.Source) (err error) {
	if err = r.Bundle.ReadSrc(src); err == nil {
		r.App = r.Bundle.App
		r.Func = r.Start
	}
	return
}

func (r *htmlBundleRunner) Start(args ...string) (err error) {
	return r.api.StartHtml(r.Bundle.Package, strings.Join(args, " "))
}
