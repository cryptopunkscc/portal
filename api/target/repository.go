package target

type Repository interface{ Get(src string) []Source }

type Repositories []Repository

var _ Repository = Repositories{}

func (r Repositories) Get(src string) (sources []Source) {
	for _, repository := range r {
		if sources = repository.Get(src); len(sources) != 0 {
			return
		}
	}
	return
}

func (c *Cache[T]) Repository() *CacheRepository[T] {
	return &CacheRepository[T]{Cache: c}
}

type CacheRepository[T Portal_] struct{ *Cache[T] }

var _ Repository = &CacheRepository[Portal_]{}

func (r *CacheRepository[T]) Get(src string) (out []Source) {
	if portal, ok := r.Cache.Get(src); ok {
		out = append(out, portal)
	}
	return
}

type SourcesRepository[T Portal_] struct {
	Sources []Source
	Resolve[T]
}

var _ Repository = &SourcesRepository[Portal_]{}

func (r *SourcesRepository[T]) Get(src string) (out []Source) {
	for _, portal := range r.Resolve.List(r.Sources...) {
		if portal.Manifest().Match(src) {
			out = append(out, portal)
		}
	}
	return
}

func (r *SourcesRepository[T]) First(src string) (out T) {
	list := r.Resolve.List(r.Sources...)
	for _, out = range list {
		if out.Manifest().Match(src) {
			return
		}
	}
	return
}

type FileRepository struct{ File }

var _ Repository = &FileRepository{}

func (r *FileRepository) Get(src string) (out []Source) {
	if file, err := r.File(src); err == nil {
		out = append(out, file)
	}
	return
}
