package install

import "github.com/cryptopunkscc/portal/api/target"

func (i Runner) generateTokenFor(bundle target.Portal_) (err error) {
	_, err = i.Tokens.Resolve(bundle.Manifest().Package)
	return
}
