package golang

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestWatcher_Run(t *testing.T) {
	t.SkipNow()
	src, _ := os.Getwd()
	wd, _ := FindProjectRoot(src)
	target := filepath.Join(wd, "pkg/go")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events, err := NewWatcher().Run(ctx, target)
	if err != nil {
		t.Fatal(err)
	}
	for event := range events {
		log.Println(event)
	}
}
