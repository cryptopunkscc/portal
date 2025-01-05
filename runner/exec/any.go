package exec

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
)

func Any(runner func(string) string) target.Runner[target.Portal_] {
	return Runner[target.Portal_](func(portal target.Portal_) ([]string, error) {
		schema := portal.Manifest().Schema
		r := runner(schema)
		if r == "" {
			return nil, fmt.Errorf("unknown schema %v", schema)
		}
		return []string{r, "o", portal.Abs()}, nil
	})
}
