package vcs

import "runtime/debug"

func ReadBuildInfo() (i BuildInfo) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				i.Revision = setting.Value
			case "vcs.modified":
				i.Modified = setting.Value
			case "vcs.time":
				i.Time = setting.Value
			}
		}
	}
	return
}

type BuildInfo struct {
	Revision string
	Time     string
	Modified string
}
