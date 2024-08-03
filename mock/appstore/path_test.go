package appstore

import "testing"

func TestPath(t *testing.T) {
	path, err := Path("launcher")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}
