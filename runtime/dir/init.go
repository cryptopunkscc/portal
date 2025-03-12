package dir

func Init(home string) {
	Bin = mk(home, "bin")
	App = mk(home, "app")
	Token = mk(home, "token")
	AppSource = src(App)
}
