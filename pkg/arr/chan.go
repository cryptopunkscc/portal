package arr

func FromChan[T any](c <-chan T) (out []T) {
	for t := range c {
		out = append(out, t)
	}
	return
}
