package target

type Provider[T Portal_] struct {
	Priority
	Repository
	Resolve[T]
}

func (r Provider[T]) Provide(src string) (out Portals[T]) {
	out = r.All(src)
	out.SortBy(r.Priority)
	out = out.Reduced()
	return
}

func (r Provider[T]) All(src string) (out Portals[T]) {
	sources := r.Repository.Get(src)
	return r.Resolve.List(sources...)
}
