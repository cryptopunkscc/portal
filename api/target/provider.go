package target

type Provider[T Portal_] struct {
	Priority
	Repository
	Resolve[T]
}

func (r Provider[T]) Provide(src string) (out Portals[T]) {
	sources := r.Repository.Get(src)
	out = r.Resolve.List(sources...)
	out.SortBy(r.Priority)
	out = out.Reduced()
	return
}
