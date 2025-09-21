package core

import (
	"errors"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/source"
)

// Deprecated
func (m *service) Install() (err error) {
	defer plog.TraceErr(&err)
	var errs []error
	i := m.Installer()
	for _, b := range bundle.Resolve_.List(
		source.Embed(apps.LauncherFS),
		source.Embed(apps.ProfileFS),
	) {
		if err = m.SetupToken(b); err != nil {
			errs = append(errs, err)
		} else if err = i.Install(b); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
