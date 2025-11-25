package env

import (
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/source"
)

type Key string

const (
	AstraldHome       Key = "ASTRALD_HOME"
	AstraldDb         Key = "ASTRALD_DB_DIR"
	ApphostAddr       Key = "APPHOST_ADDR"
	PortaldHome       Key = "PORTALD_HOME"
	PortaldTokens     Key = "PORTALD_TOKENS_DIR"
	PortaldApps       Key = "PORTALD_APPS_DIR"
	PortaldBin        Key = "PORTALD_BIN_DIR"
	PortaldConfigPath Key = "PORTALD_CONFIG_PATH"
)

func (k Key) Default(get func() string) {
	if !k.Exist() {
		v := get()
		k.Set(v)
	}
}

func (k Key) Exist() bool {
	return len(k.Get()) > 0
}

func (k Key) Get() string {
	return os.Getenv(string(k))
}

func (k Key) Set(v string) {
	err := os.Setenv(string(k), v)
	if err != nil {
		panic(err)
	}
}

func (k Key) Unset() {
	if err := os.Unsetenv(string(k)); err != nil {
		panic(err)
	}
}

func (k Key) SetDir(dir string, path ...string) {
	dir = filepath.Join(dir, filepath.Join(path...))
	k.Set(dir)
}

func (k Key) MkdirAll() (dir string) {
	abs := k.Get()
	if abs == "" {
		panic("environment variable " + k + " not set")
	}
	err := os.MkdirAll(abs, 0755)
	if err != nil {
		panic(err)
	}
	return abs
}

func (k Key) Source() target.Source {
	return source.Dir(k.MkdirAll())
}
