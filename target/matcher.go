package target

type Matcher func(arg Portal_) (ok bool)

func Match[T any](arg Portal_) (ok bool) {
	_, ok = arg.(T)
	return
}

type Priority []Matcher

func (priority Priority) Get(app Portal_) int {
	for i, match := range priority {
		if match(app) {
			return i
		}
	}
	return len(priority)
}
