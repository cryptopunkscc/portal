package apps

import "github.com/cryptopunkscc/portal/mock/appstore"

func Install(src string) error {
	return appstore.Install(src)
}
