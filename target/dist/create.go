package dist

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/source"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type CreateOpts struct {
	source.CreateOpts
	manifest.Dist
	PortalYml string
}

func Create(opts CreateOpts) (err error) {
	defer plog.TraceErr(&err)
	if opts.Path == "" {
		opts.Path = opts.Name
	}
	if opts.Name == "" {
		opts.Name = filepath.Base(opts.Path)
	}
	if err = SetDefaults(&opts.Dist); err != nil {
		return
	}
	if err = source.Create(opts.CreateOpts); err != nil {
		return
	}
	if err = os.MkdirAll(opts.Path, 0755); err != nil {
		return
	}
	b, err := yaml.Marshal(opts.Dist)
	if err != nil {
		return
	}
	if opts.PortalYml == "" {
		opts.PortalYml = "portal.yml"
	}
	portalYml := filepath.Join(opts.Path, opts.PortalYml)
	if err = os.WriteFile(portalYml, b, 0755); err != nil {
		return
	}

	return
}

func SetDefaults(dist *manifest.Dist) (err error) {
	defer plog.TraceErr(&err)

	if dist.Name == "" {
		return errors.New("missing name")
	}
	if dist.Title == "" {
		dist.Title = cases.Title(language.English).String(dist.Name)
	}
	if dist.Description == "" {
		dist.Description = "missing description"
	}
	if dist.Package == "" {
		dist.Package = strings.Join([]string{"my", "app", dist.Name}, ".")
	}
	return nil
}
