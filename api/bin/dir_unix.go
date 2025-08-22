//go:build unix

package bin

import (
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"os"
	"path/filepath"
)

func Dir() string {
	if fs2.CanWriteToDir("/bin") {
		return "/bin"
	}

	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".local/bin")
	}

	panic(err)
}
