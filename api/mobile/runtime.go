package mobile

type Runtime interface {
	Install() error
	Start()
	Stop()
	Apphost() Apphost
	App(pkg string) App
}
