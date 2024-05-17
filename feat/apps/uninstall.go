package apps

import "github.com/cryptopunkscc/go-astral-js/mock/appstore"

func Uninstall(id string) error {
	return appstore.Uninstall(id)
}
