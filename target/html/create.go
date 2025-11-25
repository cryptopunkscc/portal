package html

import (
	"context"
	"embed"
	"io/fs"

	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/template"
)

//go:embed tmpl
var templatesFS embed.FS

func Create(opts npm.CreateOpts) (err error) {
	opts.Runtime = "html"
	opts.TemplatesFS, _ = fs.Sub(templatesFS, "tmpl")
	return npm.Create(opts)
}

func RunCreate(_ context.Context, opts source.TemplateOpts, args ...string) (err error) {
	o := npm.CreateOpts{}
	o.Template = opts.Template
	if len(args) > 0 {
		o.Path = args[0]
	}
	return Create(o)
}

func ListTemplates() template.Templates {
	return template.ListFrom(templatesFS)
}
