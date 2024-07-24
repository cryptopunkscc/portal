package exec

var _env []string

func AddEnv(env ...string) {
	_env = append(_env, env...)
}

func SetEnv(env ...string) {
	_env = env
}
