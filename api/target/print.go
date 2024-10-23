package target

import "fmt"

func Sprint(source Portal_) string {
	return fmt.Sprintf("%T %s %s", source, source.Manifest().Package, source.Abs())
}
