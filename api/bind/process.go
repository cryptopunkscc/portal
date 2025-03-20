package bind

const (
	Log   = "_log"
	Sleep = "_sleep"
	Exit  = "_exit"
)

type Process interface {
	Log(arg string)
	Sleep(duration int64)
	Exit(code int)
}
