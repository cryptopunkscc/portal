package bind

type Runtime interface {
	Sys
	Apphost
}

type Module struct {
	Sys
	Apphost
}
