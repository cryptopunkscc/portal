package main

import (
	"fmt"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"os"
	"strings"
)

func main() {
	imports, err := golang.ListImports(os.Args[1])
	if err != nil {
		panic(err)
	}
	for i, s := range imports {
		print(fmt.Sprintf("%d %s (%d) ", i, s.Import, len(s.Refs)))
		//if len(s.Refs) < 2 {
		//	println(fmt.Sprintf("\t %v", s.Refs))
		//} else {
		//}
		println(":\n\t" + strings.Join(s.Refs, "\n\t"))
	}
}
