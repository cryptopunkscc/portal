package install

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
)

func (i Install) All(source target.Source, c chan<- Result) {
	defer close(c)
	for id, bundle := range apps.Resolver[target.Bundle_]().List(source) {
		err := i.Bundle(bundle)
		c <- Result{
			Id:       id,
			Error:    err,
			Manifest: *bundle.Manifest(),
		}
	}
	return
}
