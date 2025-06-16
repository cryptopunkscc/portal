package apphost

func Target(target ...string) string {
	if len(target) > 0 {
		return target[0]
	}
	return "localnode"
}
