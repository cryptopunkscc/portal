package main

import (
	"context"
)

func (a Application) Close(ctx context.Context) error {
	return a.portaldCli(ctx, "close")
}
