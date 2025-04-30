package js

import "github.com/cryptopunkscc/portal/api/target"

type Source struct{ target.Source }

var _ target.Js = Source{}

func (h Source) MainJs() {}
