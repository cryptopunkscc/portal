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
