package target

import "context"

type File func(path ...string) (source Source, err error)

func (f File) NewRun(of ...Resolve[Source]) Run[string] {
	r := Any[Runnable](of...)
	return func(ctx context.Context, src string, args ...string) (err error) {
		source, err := f(src)
		if err != nil {
			return
		}
		rr, err := r.Resolve(source)
		if err != nil {
			return
		}
		return rr.Run(ctx, args...)
	}
}
