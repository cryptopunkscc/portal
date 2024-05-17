package apps

import "github.com/cryptopunkscc/go-astral-js/mock/appstore"

func Install(src string) error {
	return appstore.Install(src)
}
