package exec

import (
	"github.com/cryptopunkscc/portal/api/target"
)

func Portal[T target.Portal_](command ...string) target.Runner[T] {
	return Runner[T](func(T) ([]string, error) { return command, nil })
}
