package backend_dev

type Backend interface {
	Run(src string) error
}
