package target

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"io/fs"
	"runtime"
	"strings"
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

type builds struct {
	Builds `json:"build,omitempty" yaml:"build,omitempty"`
}

type mapped struct {
	*Build `json:"build,omitempty" yaml:"build,omitempty"`
	nested map[string]mapped
}

func LoadBuilds(source Source) (out Builds) {
	m := mapped{}
	m.loadBuild(source)
	m.loadBuilds(source.FS())
	out = m.flatten()
	return
}

func (m *mapped) loadBuild(source Source) {
	if err := all.Unmarshalers.Load(m, source.FS(), BuildFilename); err != nil {
		return
	}
}

func (m *mapped) loadBuilds(source fs.FS) {
	bs := &builds{}
	if err := all.Unmarshalers.Load(bs, source, BuildFilename); err != nil {
		return
	}
	if d, ok := bs.Builds["default"]; ok {
		delete(bs.Builds, "default")
		m.Build = &d
	}
	m.nested = map[string]mapped{}
	for s, b := range bs.Builds {
		split := strings.Split(s, separator)
		os := split[0]
		arch := split[1:]
		mm := mapped{}
		if len(arch) == 0 {
			mm.Build = &b
		} else {
			mm.nested = map[string]mapped{}
			for _, a := range arch {
				mm.nested[a] = mapped{Build: &b}
			}
			if c, ok := m.nested[os]; ok {
				mm.Build = c.Build
			}
		}
		m.nested[os] = mm
	}
}

func (m *mapped) flatten(name ...string) (out Builds) {
	out = Builds{}
	if m.Build == nil {
		m.Build = &Build{}
	}
	if m.nested == nil {
		out[strings.Join(name, separator)] = *m.Build
		return
	}
	for s, mm := range m.nested {
		mm.Build = m.Build.merge(mm.Build)
		for ss, b := range mm.flatten(append(name, s)...) {
			out[ss] = b
		}
	}
	return
}

var separator = "-"

func MergeBuilds(build ...Build) Build {
	acc := build[0]
	list := build[1:]
	for _, next := range list {
		acc = *acc.merge(&next)
	}
	return acc
}

func (b Build) merge(build *Build) *Build {
	if build == nil {
		return &b
	}
	b.Deps = append(b.Deps, build.Deps...)
	b.Env = append(b.Env, build.Env...)
	if build.Cmd != "" {
		b.Cmd = build.Cmd
	}
	if build.Out != "" {
		b.Out = build.Out
	}
	if build.Exec != "" {
		b.Exec = build.Exec
	}
	return &b
}

func GetBuild(project Project_) (out Build) {
	ok := false
	if out, ok = project.Build()[runtime.GOOS]; !ok {
		out = project.Build()["default"]
	}
	return
}
