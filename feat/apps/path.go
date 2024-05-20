package apps

import "github.com/cryptopunkscc/go-astral-js/mock/appstore"

func Path(app string) (src string, err error) {
	return appstore.Path(app)
}
