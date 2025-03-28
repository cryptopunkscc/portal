package config

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	plog.Verbosity = 100
	dir := initTestDir(t)
	type testCase struct {
		initial Loader[any]
		name    string
		path    []string
		expect  Loader[any]
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "not found by non existing abs file path",
			path:    []string{dir, "not-existing-config"},
			initial: Loader[any]{},
			expect:  Loader[any]{},
			wantErr: true,
		},
		{
			name:    "not found by non existing local file path",
			path:    []string{"not-existing-config"},
			initial: Loader[any]{},
			expect: Loader[any]{
				wd: filepath.Dir(dir),
			},
			wantErr: true,
		},
		{
			name: "not found in working dir",
			initial: Loader[any]{
				File: "not-existing-config",
			},
			expect: Loader[any]{
				wd:   filepath.Dir(dir),
				File: "not-existing-config",
				Dir:  "/",
			},
			wantErr: true,
		},
		{
			name: "not found in given dir",
			path: []string{dir, "foo", "bar", "baz"},
			initial: Loader[any]{
				File: "not-existing-config",
			},
			expect: Loader[any]{
				File: "not-existing-config",
				Dir:  "/",
			},
			wantErr: true,
		},
		{
			name: "find in working dir",
			initial: Loader[any]{
				wd:        dir,
				Unmarshal: yaml.Unmarshal,
				File:      "config",
			},
			expect: Loader[any]{
				wd:        dir,
				Unmarshal: yaml.Unmarshal,
				File:      "config",
				Dir:       dir,
			},
		},
		{
			name: "find by local file path",
			initial: Loader[any]{
				wd:        dir,
				Unmarshal: yaml.Unmarshal,
			},
			path: []string{"config"},
			expect: Loader[any]{
				wd:        dir,
				Unmarshal: yaml.Unmarshal,
				File:      "config",
				Dir:       dir,
			},
		},
		{
			name: "find by absolute dir path",
			initial: Loader[any]{
				Unmarshal: yaml.Unmarshal,
				File:      "config",
			},
			path: []string{dir},
			expect: Loader[any]{
				Unmarshal: yaml.Unmarshal,
				File:      "config",
				Dir:       dir,
			},
		},
		{
			name: "find by absolute file path",
			path: []string{dir, "config"},
			initial: Loader[any]{
				Unmarshal: yaml.Unmarshal,
			},
			expect: Loader[any]{
				Unmarshal: yaml.Unmarshal,
				File:      "config",
				Dir:       dir,
			},
		},
		{
			name: "find upward absolute dir path",
			path: []string{dir, "foo", "bar", "baz"},
			initial: Loader[any]{
				Unmarshal: yaml.Unmarshal,
				File:      "config",
			},
			expect: Loader[any]{
				Unmarshal: yaml.Unmarshal,
				File:      "config",
				Dir:       dir,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.initial.Load(tt.path...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				plog.Println(err)
			} else {
				assert.EqualValues(t, fmt.Sprint(tt.expect), fmt.Sprint(tt.initial))
			}
		})
	}
}

func initTestDir(t *testing.T) (dir string) {
	dir = test.Mkdir(t)
	if err := os.WriteFile(filepath.Join(dir, "config"), []byte{}, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "foo", "bar", "baz"), 0755); err != nil {
		panic(err)
	}
	return
}
