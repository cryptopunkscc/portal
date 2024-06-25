package apps

import "github.com/cryptopunkscc/portal/mock/appstore"

func Path(app string) (src string, err error) {
	return appstore.Path(app)
}
