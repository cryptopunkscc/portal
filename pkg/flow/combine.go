package flow

func Combine[T any](c ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, next := range c {
			for t := range next {
				out <- t
			}
		}
	}()
	return out
}

func Emit[T any](arr []T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, t := range arr {
			out <- t
		}
	}()
	return out
}
