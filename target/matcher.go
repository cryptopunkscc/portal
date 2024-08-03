package target

type Matcher func(arg Base) (ok bool)

func Match[T any](arg Base) (ok bool) {
	_, ok = arg.(T)
	return
}

type Priority []Matcher

func (priority Priority) Get(app Base) int {
	for i, match := range priority {
		if match(app) {
			return i
		}
	}
	return len(priority)
}
