package desktop

import (
	"os"
	"path/filepath"
)

func workingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func userCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return dir
}

func userConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return dir
}

func localShareDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, ".local", "share")
}
