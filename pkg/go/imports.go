package golang

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func filterByModule(imports []string, module string) (filtered []string) {
	for _, s := range imports {
		if trim := strings.TrimPrefix(s, module); trim != s {
			filtered = append(filtered, trim)
		}
	}
	return
}

func withProjectPath(imports []string, abs string) (modified []string) {
	for _, s := range imports {
		modified = append(modified, filepath.Join(abs, s))
	}
	return
}

func FindProjectRoot(abs ...string) (string, error) {
	if p, err := ResolveProject(abs...); err != nil {
		return "", err
	} else {
		return p.Dir, nil
	}
}

func sameImports(i1, i2 []string) bool {
	if len(i1) != len(i2) {
		return false
	}
	m := make(map[string]any)
	for _, i := range i1 {
		m[i] = i
	}
	for _, i := range i2 {
		if _, ok := m[i]; !ok {
			return false
		}
	}
	return true
}

func Imports(absSrc string) (imports []string, err error) {
	if filepath.Ext(absSrc) != ".go" {
		err = fmt.Errorf("%s is not a .go file", absSrc)
		return
	}

	file, err := os.Open(absSrc)
	if err != nil {
		err = plog.Err(err)
		return
	}
	defer file.Close()

	appendNext := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == ")" {
			appendNext = false
			continue
		}
		if appendNext {
			imports = append(imports, trimImport(text))
			continue
		}
		if strings.HasPrefix(text, "import (") {
			appendNext = true
			continue
		}

		src := strings.TrimPrefix(text, "import ")
		if src != text {
			imports = append(imports, trimImport(src))
			continue
		}

		if len(imports) == 0 {
			continue
		}
		return
	}
	return
}

func trimImport(s string) string {
	c := strings.Split(s, " ")
	s = c[len(c)-1]
	s = strings.Trim(s, "\"")
	return s
}
