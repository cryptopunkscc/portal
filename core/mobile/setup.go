package core

import (
	"errors"

	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/app"
	"github.com/cryptopunkscc/portal/target/source"
)

func (m *service) installApps() (err error) {
	defer plog.TraceErr(&err)
	var errs []error
	i := m.Installer()
	for _, b := range app.Resolve_.List(
		source.Embed(apps.LauncherFS),
		source.Embed(apps.ProfileFS),
		source.Embed(apps.ClaimFS),
	) {
		if err = i.Install(b); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
