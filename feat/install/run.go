package install

import "github.com/cryptopunkscc/go-astral-js/pkg/appstore"

func Run(src string) error {
	return appstore.Install(src)
}
