package bind

const (
	Log   = "_log"
	Sleep = "_sleep"
)

type Sys interface {
	Log(arg any)
	Sleep(duration int64)
}
