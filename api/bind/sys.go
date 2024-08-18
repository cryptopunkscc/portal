package bind

const (
	Log   = "_log"
	Sleep = "_sleep"
)

type Sys interface {
	Log(str string)
	Sleep(duration int64)
}
