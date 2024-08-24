package require

func NoErr[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
