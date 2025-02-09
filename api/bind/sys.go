package bind

const (
	Log   = "_log"
	Sleep = "_sleep"
	Exit  = "_exit"
)

type Sys interface {
	Log(arg any)
	Sleep(duration int64)
	Exit(code int)
}
