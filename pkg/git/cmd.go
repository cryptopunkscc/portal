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
	reader, writer := io.Pipe()
	cmd.Env = append(os.Environ(), "TZ=UTC")
	cmd.Stdout = writer
	cmd.Stderr = os.Stderr
	scanner := bufio.NewScanner(reader)
	if err = cmd.Start(); err != nil {
		return
	}
	if !scanner.Scan() {
		err = scanner.Err()
		return
	}
	out = strings.Trim(scanner.Text(), `"`)
	_ = cmd.Wait()
	return
}
