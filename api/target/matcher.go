package target

type Matcher func(arg Portal_) (ok bool)

func Match[T any](arg Portal_) (ok bool) {
	_, ok = arg.(T)
	return
}

type Priority []Matcher

func (p Priority) Get(app Portal_) int {
	for i, match := range p {
		if match(app) {
			return i
		}
	}
	return len(p)
}

func (p Priority) Sort(order []int) (out Priority) {
	if len(order) != len(p) {
		return p
	}
	for _, i := range order {
		out = append(out, p[i])
	}
	return
}
