package list

func Chan[T any](c <-chan T) (out []T) {
	for t := range c {
		out = append(out, t)
	}
	return
}
