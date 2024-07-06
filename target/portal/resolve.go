package portal

import (
	"errors"
	"github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/source"
	"strings"
)

type Resolver[T target.Portal] struct {
	resolve target.Resolve[T]
	source  target.Source
	prefix  string
}

func NewResolver[T target.Portal](
	resolve target.Resolve[T],
	source target.Source,
) *Resolver[T] {
	return &Resolver[T]{
		resolve: resolve,
		source:  source,
		prefix:  port.PrefixStr(),
	}
}

func (f Resolver[T]) Portal(port string) (t T, err error) {
	port = strings.TrimPrefix(port, f.prefix)
	port = strings.TrimPrefix(port, ".")
	for _, t = range source.List[T](f.resolve, f.source) {
		if t.Manifest().Match(port) {
			return
		}
	}
	err = errors.New("not found")
	return
}

func (f Resolver[T]) Path(port string) (p string, err error) {
	t, err := f.Portal(port)
	if err != nil {
		return
	}
	p = t.Abs()
	return
}
