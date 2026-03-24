package mobile

type Core interface {
	Start()
	Stop()
	App(pkg string) (App, error)
}
