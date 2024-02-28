package runner

import "io"

type Bindings func() io.Closer
