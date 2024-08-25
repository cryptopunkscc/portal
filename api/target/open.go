package target

import "context"

type Open func(ctx context.Context, query string) (packages []string, err error)
