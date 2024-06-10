package port

import (
	"os"
	"strings"
)

const prefixEnv = "PORT_PREFIX"

var prefix []string

func Prefix() Port {
	return prefix
}

func PrefixStr() string {
	return Prefix().String()
}

func init() {
	p := os.Getenv(prefixEnv)
	if p != "" {
		prefix = strings.Split(p, ":")
	}
}

func InitPrefix(chunks ...string) {
	prefix = chunks
	_ = os.Setenv(prefixEnv, Prefix().String())
}
