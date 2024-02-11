package backend

type Backend interface {
	Run(path string) error
}
