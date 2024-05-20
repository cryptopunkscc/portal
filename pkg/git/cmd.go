package git

import (
	"bufio"
	"io"
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
	return ReadString(cmd)
}

func UserName(global bool) (str string) {
	cmd := exec.Command("git", "config", "user.name")
	if global {
		cmd = exec.Command("git", "config", "--global", "user.name")
	}
	str, _ = ReadString(cmd)
	return
}

func UserEmail(global bool) (str string) {
	cmd := exec.Command("git", "config", "user.email")
	if global {
		cmd = exec.Command("git", "config", "--global", "user.email")
	}
	str, _ = ReadString(cmd)
	return
}

func ReadString(cmd *exec.Cmd) (str string, err error) {
	cmd.Stderr = os.Stderr

	reader, writer := io.Pipe()
	cmd.Stdout = writer
	scanner := bufio.NewScanner(reader)
	if err = cmd.Start(); err != nil {
		return
	}
	if !scanner.Scan() {
		err = scanner.Err()
		return
	}
	str = strings.Trim(scanner.Text(), `"`)
	err = cmd.Wait()
	return
}
