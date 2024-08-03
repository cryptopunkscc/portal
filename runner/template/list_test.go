package template

import "testing"

func TestList(t *testing.T) {
	if err := List(); err != nil {
		t.Fatal(err)
	}
}
