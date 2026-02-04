package os

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadJson[T any](abs ...string) (t T, err error) {
	n := filepath.Join(abs...) + ".json"
	f, err := os.Open(n)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&t)
	return
}

func WriteJson[T any](t T, abs ...string) (err error) {
	n := filepath.Join(abs...) + ".json"
	f, err := os.Create(n)
	if err != nil {
		return
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(t)
}
