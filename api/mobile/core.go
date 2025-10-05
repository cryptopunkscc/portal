package mobile

type Core interface {
	Start()
	Ping()
	Stop()
	App(pkg string) App
}
