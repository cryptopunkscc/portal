package target

import (
	"slices"
	"strings"
)

func Op(args *[]string, arg string) (ok bool) {
	a := *args
	for i, s := range a {
		if s == arg || strings.HasSuffix(arg, "=") && strings.HasPrefix(s, arg) {
			a = slices.Delete(a, i, i+1)
			*args = a
			ok = true
			break
		}
	}
	return
}

func OpVal(args *[]string, arg string) (out string) {
	a := *args
	for i, s := range a {
		if strings.HasSuffix(arg, "=") && strings.HasPrefix(s, arg) {
			a = slices.Delete(a, i, i+1)
			*args = a
			out = strings.TrimPrefix(s, arg)
			break
		}
	}
	return
}
