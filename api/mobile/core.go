package mobile

type Core interface {
	Install() error
	Start()
	Stop()
	App(pkg string) App
}
