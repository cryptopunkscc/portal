package exec

import "strings"

var _env []string

func AddEnv(env ...string) {
	_env = append(_env, env...)
}

func SetEnv(env ...string) {
	_env = env
}

func GetEnv(key string) string {
	for _, s := range _env {
		split := strings.Split(s, "=")
		if split[0] == key {
			return split[1]
		}
	}
	return ""
}
