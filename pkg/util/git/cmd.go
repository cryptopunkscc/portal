package git

import (
	"os"
	"os/exec"
	"strings"
)

func TimestampHash() (out string, err error) {
	cmd := exec.Command("git", "--no-pager", "show",
		"--quiet",
		`--abbrev=12`,
		`--date=format-local:%Y%m%d%H%M%S`,
		`--format="%cd-%h"`,
	)
	cmd.Env = append(os.Environ(), "TZ=UTC")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	out = strings.TrimSpace(string(output))
	return
}

func UserName(global bool) (str string) {
	cmd := exec.Command("git", "config", "user.name")
	if global {
		cmd = exec.Command("git", "config", "--global", "user.name")
	}
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output))
}

func UserEmail(global bool) (str string) {
	cmd := exec.Command("git", "config", "user.email")
	if global {
		cmd = exec.Command("git", "config", "--global", "user.email")
	}
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output))
}
