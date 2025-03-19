package portald

import "context"

func Join(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
