package golang

import (
	"bufio"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"path"
	"strings"
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
		modified = append(modified, path.Join(abs, s))
	}
	return
}

func FindProjectRoot(abs string) (s string, err error) {
	if abs == "." || abs == "" || abs == "/" {
		err = plog.Errorf("cannot find root")
		return
	}
	if _, err = os.Stat(path.Join(abs, "go.mod")); err == nil {
		s = abs
		return
	}
	dir := path.Dir(abs)
	return FindProjectRoot(dir)
}

func GetModuleRoot(root string) (s string, err error) {
	goMod := path.Join(root, "go.mod")
	file, err := os.Open(goMod)
	if err != nil {
		err = plog.Err(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return "", plog.Errorf("cannot find module name")
	}

	line := strings.TrimSpace(scanner.Text())
	s = strings.TrimPrefix(line, "module ")

	if s == line {
		return "", plog.Errorf("cannot find module name")
	}

	return
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
	if path.Ext(absSrc) != ".go" {
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
