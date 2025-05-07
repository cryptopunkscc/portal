package target

import "slices"

func Op(args *[]string, arg string) (ok bool) {
	a := *args
	for i, s := range a {
		if s == arg {
			a = slices.Delete(a, i, i+1)
			*args = a
			ok = true
			break
		}
	}
	return
}
