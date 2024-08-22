package mobile

type Runtime interface {
	Start()
	Stop()
	Apphost() Apphost
	App(pkg string) App
}
