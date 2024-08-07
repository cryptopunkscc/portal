package test

import (
	"github.com/cryptopunkscc/portal/pkg/zip"
	"testing"
)

func TestZipData(t *testing.T) {
	if err := zip.Pack("data", "data.zip"); err != nil {
		t.Fatal(err)
	}
}
