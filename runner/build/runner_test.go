package build

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTest(t *testing.T) {
	assert.True(t, strings.Contains("asd", ""))
}

func Test_runner_loadTargets(t *testing.T) {
	tests := []struct {
		name    string
		targets map[string][]string
		args    []string
		wantErr bool
	}{
		{
			name:    "empty",
			args:    []string{},
			targets: map[string][]string{},
		},
		{
			name: "os only",
			args: []string{"os1", "os2"},
			targets: map[string][]string{
				"os1": nil,
				"os2": nil,
			},
		},
		{
			name: "os and arch",
			args: []string{"os1_arch", "os2_arch"},
			targets: map[string][]string{
				"os1": {"arch"},
				"os2": {"arch"},
			},
		},
		{
			name: "mixed",
			args: []string{"os", "os_arch1", "os_arch2"},
			targets: map[string][]string{
				"os": {"arch1", "arch2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner{
				targets: tt.targets,
			}
			if err := r.loadTargets(tt.args...); (err != nil) != tt.wantErr {
				t.Errorf("loadTargets() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.Equal(t, tt.targets, r.targets)
			}
		})
	}
}
