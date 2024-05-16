package apps

import "github.com/cryptopunkscc/go-astral-js/pkg/appstore"

func Install(src string) error {
	return appstore.Install(src)
}
