package mobile

type Core interface {
	Start()
	Ping()
	Setup(alias string) error
	Stop()
	App(pkg string) App
}
