package test

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"testing"
)

func TestBuild_All(t *testing.T) {
	file, err := source.File(".")
	if err != nil {
		t.Fatal(err)
	}
	builds := target.LoadBuilds(file)
	log.Println(builds)
}
