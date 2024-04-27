package backend

type Backend interface {
	Run(src string) error
}
