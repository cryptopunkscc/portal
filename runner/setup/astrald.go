package setup

import "context"

func (r *Runner) startAstrald(ctx context.Context) (err error) {
	r.log.Println("starting astrald...")
	if err = r.Runner.Start(ctx); err != nil {
		return
	}
	r.log.Println("astrald started")
	return
}
