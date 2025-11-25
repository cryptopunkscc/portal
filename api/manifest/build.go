package manifest

import (
	"io/fs"
	"path"
	"runtime"
	"slices"

	"github.com/cryptopunkscc/portal/pkg/config"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
)

type Builds struct {
	//Build  `yaml:",inline,omitempty" json:",inline,omitempty"`
	Builds InnerBuilds `yaml:"build,omitempty" json:"build,omitempty"`
}

type InnerBuilds struct {
	Build   `yaml:",inline,omitempty" json:",inline,omitempty"`
	Default Build                  `yaml:"default,omitempty" json:"default,omitempty"`
	Builds  map[string]InnerBuilds `yaml:",inline,omitempty" json:",inline,omitempty"`
}

type Build struct {
	Out    string   `json:"out,omitempty" yaml:"out,omitempty"`
	Deps   []string `json:"deps,omitempty" yaml:"deps,omitempty"`
	Env    []string `json:"env,omitempty" yaml:"env,omitempty"`
	Cmd    string   `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Exec   string   `json:"exec,omitempty" yaml:"exec,omitempty"`
	Target Target   `json:"target,omitempty" yaml:"target,omitempty"`
}

func (b *Builds) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, b) }
func (b *Builds) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(b, fs, DevFilename) }

func (b Builds) Targets() (targets Targets) {
	for os, n := range b.Builds.Builds {
		if len(n.Builds) == 0 {
			targets = append(targets, []string{os})
		}
		for arch := range n.Builds {
			targets = append(targets, []string{os, arch})
		}
	}
	return
}

type Targets [][]string

func (t Targets) Flatten() (out []string) {
	out = make([]string, len(t))
	for i, s := range t {
		out[i] = path.Join(s...)
	}
	return
}

func (b Builds) Get(targets ...string) (build Build) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	if len(targets) > 0 {
		os = targets[0]
	}
	if len(targets) > 1 {
		arch = targets[1]
	}
	merge := []*Build{&b.Builds.Build, &b.Builds.Default}
	if n, ok := b.Builds.Builds[os]; ok {
		merge = append(merge, &n.Build, &n.Default)
		if n, ok := n.Builds[arch]; ok {
			merge = append(merge, &n.Build, &n.Default)
		}
	}
	slices.Reverse(merge)
	config.Merge(&build, merge...)

	build.Target.OS = os
	build.Target.Arch = arch
	build.Target.Exec = build.Out
	return
}
