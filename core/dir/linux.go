//go:build linux && !android

package dir

import (
	"os"
)

func init() {
	// portald
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	portaldDir := mk(homeDir, ".local/share/portald")
	Init(portaldDir)

	// astrald
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	Node = mk(configDir, "astrald")
}
