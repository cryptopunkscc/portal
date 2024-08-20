package mobile

type Api interface {
	NodeRoot() string
	RequestHtml(src string) error
	Event(event *Event)
}

type Event struct {
	Msg int
	Err error
}

const (
	STARTING = iota
	STARTED
	STOPPED
)
