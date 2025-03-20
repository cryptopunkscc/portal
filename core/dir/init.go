package dir

func Init(portald string) {
	Bin = mk(portald, "bin")
	App = mk(portald, "app")
	Token = mk(portald, "token")
	AppSource = src(App)
}
