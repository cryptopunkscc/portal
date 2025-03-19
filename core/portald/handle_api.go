package portald

import (
	"github.com/cryptopunkscc/portal/core/rpc/cmd"
)

func (s *Runner[T]) Api() cmd.Handlers {
	return s.portaldApi().Plus(s.appApi()...)
}

func (s *Runner[T]) portaldApi() (handlers cmd.Handlers) {
	for _, handler := range s.publicHandlers() {
		handler.Func = "portald"
		handler.Sub = nil
		handlers = append(handlers, handler)
	}
	return handlers
}

func (s *Runner[T]) appApi() (handlers cmd.Handlers) {
	for _, app := range s.ListApps(ListAppsOpts{}) {
		m := app.Manifest()
		h := cmd.Handler{
			Func: "app",
			Name: m.Name,
			Desc: m.Description,
		}
		handlers = append(handlers, h)
	}
	return
}
