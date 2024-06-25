package apps

import "github.com/cryptopunkscc/portal/mock/appstore"

func Uninstall(id string) error {
	return appstore.Uninstall(id)
}
