package apps

import "github.com/cryptopunkscc/go-astral-js/pkg/appstore"

func Uninstall(id string) error {
	return appstore.Uninstall(id)
}
