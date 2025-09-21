package portald

import (
	apphost2 "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/api/target"
)

func (s *Service) ClaimPackage(pkg string) (t *apphost2.AccessToken, err error) {
	t, err = s.Tokens().Resolve(pkg)
	if err != nil {
		return
	}
	if s.HasUser() {
		err = s.signAppContract(t.Identity.String())
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) ClaimApp(app target.App_) (err error) {
	_, err = s.ClaimPackage(app.Manifest().Package)
	return
}
