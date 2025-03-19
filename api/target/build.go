package target

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"runtime"
)

const BuildFilename = "dev.portal"

type Build struct {
	Out  string   `json:"out,omitempty" yaml:"out,omitempty"`
	Deps []string `json:"deps,omitempty" yaml:"deps,omitempty"`
	Env  []string `json:"env,omitempty" yaml:"env,omitempty"`
	Cmd  string   `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Exec string   `json:"exec,omitempty" yaml:"exec,omitempty"`
}

type Builds map[string]Build

type build struct {
	Build Build `json:"build,omitempty" yaml:"build,omitempty"`
}

type builds struct {
	Builds Builds `json:"build,omitempty" yaml:"build,omitempty"`
}

func LoadBuilds(source Source) (out Builds) {
	out = Builds{}
	b := build{}
	if err := all.Unmarshalers.Load(&b, source.FS(), BuildFilename); err == nil {
		if b.Build.Cmd != "" || len(b.Build.Deps) > 0 {
			out["default"] = b.Build
			return
		}
	}
	bs := builds{}
	if err := all.Unmarshalers.Load(&bs, source.FS(), BuildFilename); err == nil {
		out = bs.Builds
	}
	if def, ok := out["default"]; ok {
		for key, bb := range bs.Builds {
			if key != "default" {
				out[key] = def.merge(bb)
			}
		}
	}
	return
}

func (b Build) merge(build Build) Build {
	b.Deps = append(b.Deps, build.Deps...)
	b.Env = append(b.Env, build.Env...)
	if build.Cmd != "" {
		b.Cmd = build.Cmd
	}
	if build.Out != "" {
		b.Out = build.Out
	}
	return b
}

func GetBuild(project Project_) (out Build) {
	ok := false
	if out, ok = project.Build()[runtime.GOOS]; !ok {
		out = project.Build()["default"]
	}
	return
}
