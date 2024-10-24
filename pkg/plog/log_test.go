package plog

import (
	"context"
	"testing"
)

func TestApi(t *testing.T) {
	ctx := context.Background()
	log := New().W().Scope("root").Set(&ctx)
	log.Println("message 1")
	nested(ctx)
	log.F().Println("message 3")
}

func nested(ctx context.Context) {
	Get(ctx).Scope("nested").Set(&ctx)
	Get(ctx).I().Println("message 2")
}

func TestImmutability(t *testing.T) {
	ctx := context.Background()
	Get(ctx).Scope("foo").Set(&ctx)
	log1 := Get(ctx).Scope("bar")
	log2 := Get(ctx).Scope("baz")
	log3 := log2.Scope("last")
	log3.Printf("message")
	log2.Printf("message")
	log1.Printf("message")
}
