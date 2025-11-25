package gpg

import (
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/pkg/exec"
)

func Sign(path string) {
	dir, file := filepath.Split(path)
	_ = os.Remove(path + ".sig")
	_ = exec.Run(dir, "gpg",
		"--sign",
		"--detach-sign",
		"--verbose",
		"--digest-algo", "sha512",
		file,
	)
}
