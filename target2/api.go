package target2

type Resolve[T any] func(src Source) (result T, err error)
