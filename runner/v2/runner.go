package runner

import "github.com/cryptopunkscc/portal/source"

type Loader interface {
	Load(source source.Ref) []Runner
}

type Runner interface {
	Run(args ...string) error
}
