package uninstall

import "github.com/cryptopunkscc/go-astral-js/pkg/appstore"

func Run(id string) error {
	return appstore.Uninstall(id)
}
