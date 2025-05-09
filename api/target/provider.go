package target

type Provider[T Portal_] struct {
	Priority
	Repository
	Filter[T]
	Resolve[T]
}

type Filter[T Portal_] func(T) bool

func (r Provider[T]) Provide(src string) (out Portals[T]) {
	out = r.All(src)
	out.SortBy(r.Priority)
	out = r.filter(out)
	out = out.Reduced()
	return
}

func (r Provider[T]) All(src string) (out Portals[T]) {
	sources := r.Repository.Get(src)
	return r.Resolve.List(sources...)
}

func (r Provider[T]) filter(in Portals[T]) (out Portals[T]) {
	if r.Filter == nil {
		return in
	}
	for _, t := range in {
		if r.Filter(t) {
			out = append(out, t)
		}
	}
	return
}
