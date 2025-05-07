package target

import (
	"fmt"
	"strings"
)

func Sprint(source Portal_) (s string) {
	s = fmt.Sprintf("%T %s %s", source, source.Manifest().Package, source.Abs())
	s = strings.ReplaceAll(s, "github.com/cryptopunkscc/portal/api/target.", "")
	return
}
